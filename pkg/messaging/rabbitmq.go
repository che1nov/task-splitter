package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"task-splitter/config"
	"time"

	"github.com/streadway/amqp"
)

// MessagingService интерфейс для работы с сообщениями
type MessagingService interface {
	PublishSplitTaskRequest(message SplitTaskMessage) error
	ConsumeSplitTaskRequests(handler func(SplitTaskMessage) error) error
	PublishSplitTaskResponse(message SplitTaskResponse) error
	ConsumeSplitTaskResponses(handler func(SplitTaskResponse) error) error
	Close() error
}

// SplitTaskMessage представляет сообщение запроса на разбивку задачи
type SplitTaskMessage struct {
	RequestID  string `json:"request_id"`
	TaskID     uint   `json:"task_id"`
	UserID     uint   `json:"user_id"`
	Text       string `json:"text"`
	TemplateID *uint  `json:"template_id"`
}

// SplitTaskResponse представляет ответ на разбивку задачи
type SplitTaskResponse struct {
	RequestID string `json:"request_id"`
	TaskID    uint   `json:"task_id"`
	UserID    uint   `json:"user_id"`
	Status    string `json:"status"` // completed, failed
	Result    string `json:"result"`
	Error     string `json:"error"`
}

// rabbitMQService реализация MessagingService
type rabbitMQService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	cfg     config.RabbitMQConfig
}

// NewRabbitMQService создает новый сервис RabbitMQ
func NewRabbitMQService() MessagingService {
	// Загружаем конфигурацию из переменных окружения
	cfg := config.Load().RabbitMQ

	// Пытаемся подключиться с retry логикой
	conn, err := connectWithRetry(cfg, 5, 2*time.Second)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ after retries: %v", err)
		// Возвращаем заглушку для разработки
		return &mockMessagingService{}
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open channel: %v", err)
		conn.Close()
		return &mockMessagingService{}
	}

	service := &rabbitMQService{
		conn:    conn,
		channel: ch,
		cfg:     cfg,
	}

	// Объявляем очереди
	if err := service.declareQueues(); err != nil {
		log.Printf("Failed to declare queues: %v", err)
		service.Close()
		return &mockMessagingService{}
	}

	return service
}

// connectWithRetry пытается подключиться к RabbitMQ с повторными попытками
func connectWithRetry(cfg config.RabbitMQConfig, maxRetries int, delay time.Duration) (*amqp.Connection, error) {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s%s",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.VHost))
		if err == nil {
			log.Printf("Successfully connected to RabbitMQ on attempt %d", i+1)
			return conn, nil
		}

		lastErr = err
		log.Printf("Failed to connect to RabbitMQ (attempt %d/%d): %v", i+1, maxRetries, err)

		if i < maxRetries-1 {
			log.Printf("Retrying in %v...", delay)
			time.Sleep(delay)
		}
	}

	return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, lastErr)
}

// declareQueues объявляет необходимые очереди
func (s *rabbitMQService) declareQueues() error {
	// Очередь для запросов на разбивку задач
	_, err := s.channel.QueueDeclare(
		"task.split.request", // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare request queue: %w", err)
	}

	// Очередь для ответов на разбивку задач
	_, err = s.channel.QueueDeclare(
		"task.split.response", // name
		true,                  // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare response queue: %w", err)
	}

	return nil
}

func (s *rabbitMQService) PublishSplitTaskRequest(message SplitTaskMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = s.channel.Publish(
		"",                   // exchange
		"task.split.request", // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Сохраняем сообщение на диск
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published split task request: %s", message.RequestID)
	return nil
}

func (s *rabbitMQService) ConsumeSplitTaskRequests(handler func(SplitTaskMessage) error) error {
	msgs, err := s.channel.Consume(
		"task.split.request", // queue
		"",                   // consumer
		false,                // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for d := range msgs {
			var message SplitTaskMessage
			if err := json.Unmarshal(d.Body, &message); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				d.Nack(false, false)
				continue
			}

			if err := handler(message); err != nil {
				log.Printf("Failed to process message: %v", err)
				d.Nack(false, true) // Повторяем обработку
			} else {
				d.Ack(false)
			}
		}
	}()

	return nil
}

func (s *rabbitMQService) PublishSplitTaskResponse(message SplitTaskResponse) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = s.channel.Publish(
		"",                    // exchange
		"task.split.response", // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Сохраняем сообщение на диск
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published split task response: %s", message.RequestID)
	return nil
}

func (s *rabbitMQService) ConsumeSplitTaskResponses(handler func(SplitTaskResponse) error) error {
	msgs, err := s.channel.Consume(
		"task.split.response", // queue
		"",                    // consumer
		false,                 // auto-ack
		false,                 // exclusive
		false,                 // no-local
		false,                 // no-wait
		nil,                   // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for d := range msgs {
			var message SplitTaskResponse
			if err := json.Unmarshal(d.Body, &message); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				d.Nack(false, false)
				continue
			}

			if err := handler(message); err != nil {
				log.Printf("Failed to process message: %v", err)
				d.Nack(false, true) // Повторяем обработку
			} else {
				d.Ack(false)
			}
		}
	}()

	return nil
}

func (s *rabbitMQService) Close() error {
	if s.channel != nil {
		s.channel.Close()
	}
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

// mockMessagingService заглушка для разработки
type mockMessagingService struct{}

func (m *mockMessagingService) PublishSplitTaskRequest(message SplitTaskMessage) error {
	log.Printf("Mock: Publishing split task request: %s", message.RequestID)
	return nil
}

func (m *mockMessagingService) ConsumeSplitTaskRequests(handler func(SplitTaskMessage) error) error {
	log.Printf("Mock: Consuming split task requests")
	return nil
}

func (m *mockMessagingService) PublishSplitTaskResponse(message SplitTaskResponse) error {
	log.Printf("Mock: Publishing split task response: %s", message.RequestID)
	return nil
}

func (m *mockMessagingService) ConsumeSplitTaskResponses(handler func(SplitTaskResponse) error) error {
	log.Printf("Mock: Consuming split task responses")
	return nil
}

func (m *mockMessagingService) Close() error {
	return nil
}
