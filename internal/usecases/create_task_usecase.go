package usecases

import (
	"context"
	"task-splitter/internal/domain"
	"task-splitter/internal/dto"
	"task-splitter/pkg/logger"
)

// CreateTaskUseCase Use Case для создания задачи
type CreateTaskUseCase struct {
	taskRepo TaskRepository
	log      logger.Logger
}

// NewCreateTaskUseCase создает новый UseCase для создания задачи
func NewCreateTaskUseCase(taskRepo TaskRepository, log logger.Logger) *CreateTaskUseCase {
	return &CreateTaskUseCase{
		taskRepo: taskRepo,
		log:      log,
	}
}

// Execute создает новую задачу
func (uc *CreateTaskUseCase) Execute(ctx context.Context, input dto.CreateTaskInput) (dto.TaskOutput, error) {
	uc.log.InfoContext(ctx, "Создание задачи", "user_id", input.UserID, "title", input.Title)

	// Создаем доменную сущность с валидацией
	task, err := domain.NewTask(input.UserID, input.Title, input.Description, input.Priority, input.DueDate)
	if err != nil {
		uc.log.ErrorContext(ctx, "Ошибка валидации задачи", "error", err)
		return dto.TaskOutput{}, err
	}

	// Сохраняем в БД
	if err := uc.taskRepo.CreateTask(ctx, task); err != nil {
		uc.log.ErrorContext(ctx, "Ошибка сохранения задачи в БД", "error", err)
		return dto.TaskOutput{}, err
	}

	uc.log.InfoContext(ctx, "Задача успешно создана", "task_id", task.ID, "user_id", task.UserID)

	return dto.TaskOutput{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		DueDate:     task.DueDate,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}, nil
}

