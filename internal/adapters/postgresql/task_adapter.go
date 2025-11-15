package postgresql

import (
	"context"
	"task-splitter/internal/domain"
	"task-splitter/pkg/logger"

	"gorm.io/gorm"
)

// TaskAdapter адаптер для работы с задачами в PostgreSQL
type TaskAdapter struct {
	db  *gorm.DB
	log logger.Logger
}

// NewTaskAdapter создает новый адаптер задач
func NewTaskAdapter(db *gorm.DB, log logger.Logger) *TaskAdapter {
	return &TaskAdapter{
		db:  db,
		log: log,
	}
}

// CreateTask создает задачу в БД
func (a *TaskAdapter) CreateTask(ctx context.Context, task domain.Task) error {
	model := TaskModel{
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		DueDate:     task.DueDate,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}

	if err := a.db.WithContext(ctx).Create(&model).Error; err != nil {
		a.log.ErrorContext(ctx, "Ошибка создания задачи в БД", "error", err)
		return err
	}

	return nil
}

// GetTaskByID получает задачу по ID
func (a *TaskAdapter) GetTaskByID(ctx context.Context, id uint) (domain.Task, error) {
	var model TaskModel
	if err := a.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Task{}, domain.ErrTaskNotFound
		}
		a.log.ErrorContext(ctx, "Ошибка получения задачи по ID", "id", id, "error", err)
		return domain.Task{}, err
	}

	return toDomainTask(model), nil
}

// GetTasksByUserID получает задачи пользователя
func (a *TaskAdapter) GetTasksByUserID(ctx context.Context, userID uint, limit, offset int) ([]domain.Task, error) {
	var models []TaskModel
	if err := a.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		a.log.ErrorContext(ctx, "Ошибка получения задач пользователя", "user_id", userID, "error", err)
		return nil, err
	}

	tasks := make([]domain.Task, len(models))
	for i, model := range models {
		tasks[i] = toDomainTask(model)
	}

	return tasks, nil
}

// GetTaskWithSubtasks получает задачу с подзадачами
func (a *TaskAdapter) GetTaskWithSubtasks(ctx context.Context, id uint) (domain.Task, []domain.Subtask, error) {
	var taskModel TaskModel
	if err := a.db.WithContext(ctx).First(&taskModel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Task{}, nil, domain.ErrTaskNotFound
		}
		a.log.ErrorContext(ctx, "Ошибка получения задачи", "id", id, "error", err)
		return domain.Task{}, nil, err
	}

	var subtaskModels []SubtaskModel
	if err := a.db.WithContext(ctx).
		Where("task_id = ?", id).
		Order("\"order\" ASC").
		Find(&subtaskModels).Error; err != nil {
		a.log.ErrorContext(ctx, "Ошибка получения подзадач", "task_id", id, "error", err)
		return domain.Task{}, nil, err
	}

	task := toDomainTask(taskModel)
	subtasks := make([]domain.Subtask, len(subtaskModels))
	for i, model := range subtaskModels {
		subtasks[i] = toDomainSubtask(model)
	}

	return task, subtasks, nil
}

// UpdateTask обновляет задачу
func (a *TaskAdapter) UpdateTask(ctx context.Context, task domain.Task) error {
	model := TaskModel{
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

	if err := a.db.WithContext(ctx).Save(&model).Error; err != nil {
		a.log.ErrorContext(ctx, "Ошибка обновления задачи", "id", task.ID, "error", err)
		return err
	}

	return nil
}

// DeleteTask удаляет задачу
func (a *TaskAdapter) DeleteTask(ctx context.Context, id uint) error {
	if err := a.db.WithContext(ctx).Delete(&TaskModel{}, id).Error; err != nil {
		a.log.ErrorContext(ctx, "Ошибка удаления задачи", "id", id, "error", err)
		return err
	}

	return nil
}

// toDomainTask конвертирует модель БД в доменную модель
func toDomainTask(model TaskModel) domain.Task {
	return domain.Task{
		ID:          model.ID,
		UserID:      model.UserID,
		Title:       model.Title,
		Description: model.Description,
		Status:      model.Status,
		Priority:    model.Priority,
		DueDate:     model.DueDate,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

// toDomainSubtask конвертирует модель БД в доменную модель подзадачи
func toDomainSubtask(model SubtaskModel) domain.Subtask {
	return domain.Subtask{
		ID:            model.ID,
		TaskID:        model.TaskID,
		Title:         model.Title,
		Description:   model.Description,
		Status:        model.Status,
		Priority:      model.Priority,
		Order:         model.Order,
		EstimatedTime: model.EstimatedTime,
		ActualTime:    model.ActualTime,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
}

