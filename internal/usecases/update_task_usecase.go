package usecases

import (
	"context"
	"task-splitter/internal/domain"
	"task-splitter/internal/dto"
	"task-splitter/pkg/logger"
)

// UpdateTaskUseCase Use Case для обновления задачи
type UpdateTaskUseCase struct {
	taskRepo TaskRepository
	log      logger.Logger
}

// NewUpdateTaskUseCase создает новый UseCase для обновления задачи
func NewUpdateTaskUseCase(taskRepo TaskRepository, log logger.Logger) *UpdateTaskUseCase {
	return &UpdateTaskUseCase{
		taskRepo: taskRepo,
		log:      log,
	}
}

// Execute обновляет задачу
func (uc *UpdateTaskUseCase) Execute(ctx context.Context, input dto.UpdateTaskInput) (dto.TaskOutput, error) {
	uc.log.InfoContext(ctx, "Обновление задачи", "task_id", input.TaskID)

	// Получаем задачу из БД
	task, err := uc.taskRepo.GetTaskByID(ctx, input.TaskID)
	if err != nil {
		uc.log.ErrorContext(ctx, "Задача не найдена", "task_id", input.TaskID, "error", err)
		return dto.TaskOutput{}, domain.ErrTaskNotFound
	}

	// Обновляем задачу
	if err := task.Update(input.Title, input.Description, input.Status, input.Priority, input.DueDate); err != nil {
		uc.log.ErrorContext(ctx, "Ошибка валидации при обновлении", "task_id", input.TaskID, "error", err)
		return dto.TaskOutput{}, err
	}

	// Сохраняем изменения
	if err := uc.taskRepo.UpdateTask(ctx, task); err != nil {
		uc.log.ErrorContext(ctx, "Ошибка сохранения изменений", "task_id", input.TaskID, "error", err)
		return dto.TaskOutput{}, err
	}

	uc.log.InfoContext(ctx, "Задача успешно обновлена", "task_id", task.ID)

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

