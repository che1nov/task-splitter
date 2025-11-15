package usecases

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"task-splitter/internal/domain"
	"task-splitter/internal/dto"
	"task-splitter/pkg/logger"
)

// SplitTaskUseCase Use Case для разбивки задачи
type SplitTaskUseCase struct {
	taskRepo        TaskRepository
	splitRequestRepo TaskSplitRequestRepository
	messaging       MessagingAdapter
	log             logger.Logger
}

// NewSplitTaskUseCase создает новый UseCase для разбивки задачи
func NewSplitTaskUseCase(
	taskRepo TaskRepository,
	splitRequestRepo TaskSplitRequestRepository,
	messaging MessagingAdapter,
	log logger.Logger,
) *SplitTaskUseCase {
	return &SplitTaskUseCase{
		taskRepo:        taskRepo,
		splitRequestRepo: splitRequestRepo,
		messaging:       messaging,
		log:             log,
	}
}

// Execute отправляет задачу на разбивку
func (uc *SplitTaskUseCase) Execute(ctx context.Context, input dto.SplitTaskInput) (dto.SplitTaskOutput, error) {
	uc.log.InfoContext(ctx, "Разбивка задачи", "task_id", input.TaskID, "user_id", input.UserID)

	// Проверяем существование задачи и права доступа
	task, err := uc.taskRepo.GetTaskByID(ctx, input.TaskID)
	if err != nil {
		uc.log.ErrorContext(ctx, "Задача не найдена", "task_id", input.TaskID, "error", err)
		return dto.SplitTaskOutput{}, domain.ErrTaskNotFound
	}

	if !task.BelongsToUser(input.UserID) {
		uc.log.WarnContext(ctx, "Попытка разбить чужую задачу", "task_id", input.TaskID, "user_id", input.UserID)
		return dto.SplitTaskOutput{}, domain.ErrAccessDenied
	}

	// Генерируем уникальный ID запроса
	requestID, err := generateRequestID()
	if err != nil {
		uc.log.ErrorContext(ctx, "Ошибка генерации ID запроса", "error", err)
		return dto.SplitTaskOutput{}, err
	}

	// Создаем запрос на разбивку с валидацией
	splitRequest, err := domain.NewTaskSplitRequest(requestID, input.TaskID, input.UserID, input.Text, input.TemplateID)
	if err != nil {
		uc.log.ErrorContext(ctx, "Ошибка валидации запроса разбивки", "error", err)
		return dto.SplitTaskOutput{}, err
	}

	// Публикуем сообщение в очередь
	if err := uc.messaging.PublishSplitTaskRequest(ctx, requestID, input.TaskID, input.UserID, input.Text, input.TemplateID); err != nil {
		uc.log.ErrorContext(ctx, "Ошибка публикации сообщения", "error", err)
		return dto.SplitTaskOutput{}, err
	}

	// Сохраняем запрос в БД
	if err := uc.splitRequestRepo.CreateSplitRequest(ctx, splitRequest); err != nil {
		uc.log.WarnContext(ctx, "Ошибка сохранения запроса в БД (сообщение уже отправлено)", "error", err)
	}

	uc.log.InfoContext(ctx, "Запрос на разбивку создан", "request_id", requestID, "task_id", input.TaskID)

	return dto.SplitTaskOutput{
		RequestID: requestID,
		Status:    domain.SplitStatusPending,
		Message:   "Запрос на разбивку задачи успешно отправлен",
	}, nil
}

// generateRequestID генерирует уникальный ID запроса
func generateRequestID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

