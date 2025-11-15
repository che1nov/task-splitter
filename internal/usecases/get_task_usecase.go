package usecases

import (
	"context"
	"task-splitter/internal/dto"
	"task-splitter/pkg/logger"
)

// GetTaskUseCase Use Case для получения задачи
type GetTaskUseCase struct {
	taskRepo TaskRepository
	log      logger.Logger
}

// NewGetTaskUseCase создает новый UseCase для получения задачи
func NewGetTaskUseCase(taskRepo TaskRepository, log logger.Logger) *GetTaskUseCase {
	return &GetTaskUseCase{
		taskRepo: taskRepo,
		log:      log,
	}
}

// Execute получает задачу с подзадачами по ID
func (uc *GetTaskUseCase) Execute(ctx context.Context, taskID uint) (dto.TaskOutput, error) {
	uc.log.DebugContext(ctx, "Получение задачи", "task_id", taskID)

	// Получаем задачу с подзадачами из БД
	task, subtasks, err := uc.taskRepo.GetTaskWithSubtasks(ctx, taskID)
	if err != nil {
		uc.log.ErrorContext(ctx, "Ошибка получения задачи", "task_id", taskID, "error", err)
		return dto.TaskOutput{}, err
	}

	// Формируем DTO подзадач
	subtasksDTO := make([]dto.SubtaskOutput, len(subtasks))
	for i, subtask := range subtasks {
		subtasksDTO[i] = dto.SubtaskOutput{
			ID:            subtask.ID,
			TaskID:        subtask.TaskID,
			Title:         subtask.Title,
			Description:   subtask.Description,
			Status:        subtask.Status,
			Priority:      subtask.Priority,
			Order:         subtask.Order,
			EstimatedTime: subtask.EstimatedTime,
			ActualTime:    subtask.ActualTime,
			CreatedAt:     subtask.CreatedAt,
			UpdatedAt:     subtask.UpdatedAt,
		}
	}

	uc.log.DebugContext(ctx, "Задача получена", "task_id", taskID, "subtasks_count", len(subtasks))

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
		Subtasks:    subtasksDTO,
	}, nil
}

