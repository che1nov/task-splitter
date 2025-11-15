package usecases

import (
	"context"
	"task-splitter/internal/domain"
	"task-splitter/pkg/logger"
)

// DeleteTaskUseCase Use Case для удаления задачи
type DeleteTaskUseCase struct {
	taskRepo TaskRepository
	log      logger.Logger
}

// NewDeleteTaskUseCase создает новый UseCase для удаления задачи
func NewDeleteTaskUseCase(taskRepo TaskRepository, log logger.Logger) *DeleteTaskUseCase {
	return &DeleteTaskUseCase{
		taskRepo: taskRepo,
		log:      log,
	}
}

// Execute удаляет задачу
func (uc *DeleteTaskUseCase) Execute(ctx context.Context, taskID, userID uint) error {
	uc.log.InfoContext(ctx, "Удаление задачи", "task_id", taskID, "user_id", userID)

	// Получаем задачу для проверки владельца
	task, err := uc.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		uc.log.ErrorContext(ctx, "Задача не найдена", "task_id", taskID, "error", err)
		return domain.ErrTaskNotFound
	}

	// Проверяем принадлежность задачи пользователю
	if !task.BelongsToUser(userID) {
		uc.log.WarnContext(ctx, "Попытка удалить чужую задачу", "task_id", taskID, "user_id", userID, "owner_id", task.UserID)
		return domain.ErrAccessDenied
	}

	// Удаляем задачу
	if err := uc.taskRepo.DeleteTask(ctx, taskID); err != nil {
		uc.log.ErrorContext(ctx, "Ошибка удаления задачи", "task_id", taskID, "error", err)
		return err
	}

	uc.log.InfoContext(ctx, "Задача успешно удалена", "task_id", taskID)

	return nil
}

