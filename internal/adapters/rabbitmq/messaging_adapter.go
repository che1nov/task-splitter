package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"task-splitter/pkg/logger"

	"github.com/streadway/amqp"
)

// MessagingAdapter адаптер для работы с RabbitMQ
type MessagingAdapter struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	log     logger.Logger
}

// NewMessagingAdapter создает новый адаптер обмена сообщениями
func NewMessagingAdapter(amqpURL string, log logger.Logger) (*MessagingAdapter, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Error("Ошибка подключения к RabbitMQ", "error", err)
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Error("Ошибка создания канала RabbitMQ", "error", err)
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	// Объявляем очередь для задач разбивки
	_, err = channel.QueueDeclare(
		"task_split_queue", // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		log.Error("Ошибка объявления очереди", "error", err)
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	log.Info("RabbitMQ адаптер успешно инициализирован")

	return &MessagingAdapter{
		conn:    conn,
		channel: channel,
		log:     log,
	}, nil
}

// SplitTaskMessage сообщение для разбивки задачи
type SplitTaskMessage struct {
	RequestID  string `json:"request_id"`
	TaskID     uint   `json:"task_id"`
	UserID     uint   `json:"user_id"`
	Text       string `json:"text"`
	TemplateID *uint  `json:"template_id"`
}

// PublishSplitTaskRequest публикует запрос на разбивку задачи
func (a *MessagingAdapter) PublishSplitTaskRequest(ctx context.Context, requestID string, taskID, userID uint, text string, templateID *uint) error {
	message := SplitTaskMessage{
		RequestID:  requestID,
		TaskID:     taskID,
		UserID:     userID,
		Text:       text,
		TemplateID: templateID,
	}

	body, err := json.Marshal(message)
	if err != nil {
		a.log.ErrorContext(ctx, "Ошибка сериализации сообщения", "error", err)
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = a.channel.Publish(
		"",                 // exchange
		"task_split_queue", // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err != nil {
		a.log.ErrorContext(ctx, "Ошибка публикации сообщения", "error", err)
		return fmt.Errorf("failed to publish message: %w", err)
	}

	a.log.InfoContext(ctx, "Сообщение успешно опубликовано", "request_id", requestID, "task_id", taskID)
	return nil
}

// Close закрывает соединение с RabbitMQ
func (a *MessagingAdapter) Close() error {
	if a.channel != nil {
		if err := a.channel.Close(); err != nil {
			a.log.Error("Ошибка закрытия канала RabbitMQ", "error", err)
		}
	}
	if a.conn != nil {
		if err := a.conn.Close(); err != nil {
			a.log.Error("Ошибка закрытия соединения RabbitMQ", "error", err)
			return err
		}
	}

	a.log.Info("RabbitMQ соединение закрыто")
	return nil
}

