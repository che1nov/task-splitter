package usecases

import (
	"context"
	"task-splitter/internal/dto"
	"task-splitter/pkg/logger"
)

// GetTasksUseCase Use Case для получения списка задач
type GetTasksUseCase struct {
	taskRepo TaskRepository
	log      logger.Logger
}

// NewGetTasksUseCase создает новый UseCase для получения списка задач
func NewGetTasksUseCase(taskRepo TaskRepository, log logger.Logger) *GetTasksUseCase {
	return &GetTasksUseCase{
		taskRepo: taskRepo,
		log:      log,
	}
}

// Execute получает список задач пользователя
func (uc *GetTasksUseCase) Execute(ctx context.Context, input dto.GetTasksByUserInput) ([]dto.TaskOutput, error) {
	uc.log.DebugContext(ctx, "Получение списка задач", "user_id", input.UserID, "limit", input.Limit, "offset", input.Offset)

	// Получаем задачи из БД
	tasks, err := uc.taskRepo.GetTasksByUserID(ctx, input.UserID, input.Limit, input.Offset)
	if err != nil {
		uc.log.ErrorContext(ctx, "Ошибка получения задач", "user_id", input.UserID, "error", err)
		return nil, err
	}

	// Формируем DTO
	result := make([]dto.TaskOutput, len(tasks))
	for i, task := range tasks {
		result[i] = dto.TaskOutput{
			ID:          task.ID,
			UserID:      task.UserID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			Priority:    task.Priority,
			DueDate:     task.DueDate,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		}
	}

	uc.log.DebugContext(ctx, "Задачи получены", "user_id", input.UserID, "count", len(tasks))

	return result, nil
}

