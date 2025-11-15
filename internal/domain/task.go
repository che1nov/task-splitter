package domain

import (
	"strings"
	"time"
)

// Task представляет основную задачу
type Task struct {
	ID          uint
	UserID      uint
	Title       string
	Description string
	Status      string
	Priority    string
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewTask создает новую задачу с валидацией
func NewTask(userID uint, title, description, priority string, dueDate *time.Time) (Task, error) {
	// Валидация заголовка
	title = strings.TrimSpace(title)
	if title == "" {
		return Task{}, ErrInvalidTitle
	}

	// Валидация приоритета
	if priority == "" {
		priority = PriorityMedium // Приоритет по умолчанию
	}
	if !isValidPriority(priority) {
		return Task{}, ErrInvalidPriority
	}

	return Task{
		UserID:      userID,
		Title:       title,
		Description: strings.TrimSpace(description),
		Status:      StatusPending, // Статус по умолчанию
		Priority:    priority,
		DueDate:     dueDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// Update обновляет задачу
func (t *Task) Update(title, description, status, priority string, dueDate *time.Time) error {
	// Валидация заголовка
	if title != "" {
		title = strings.TrimSpace(title)
		if title == "" {
			return ErrInvalidTitle
		}
		t.Title = title
	}

	// Обновление описания
	if description != "" {
		t.Description = strings.TrimSpace(description)
	}

	// Валидация и обновление статуса
	if status != "" {
		if !isValidTaskStatus(status) {
			return ErrInvalidStatus
		}
		t.Status = status
	}

	// Валидация и обновление приоритета
	if priority != "" {
		if !isValidPriority(priority) {
			return ErrInvalidPriority
		}
		t.Priority = priority
	}

	// Обновление срока
	if dueDate != nil {
		t.DueDate = dueDate
	}

	t.UpdatedAt = time.Now()
	return nil
}

// MarkAsInProgress помечает задачу как выполняющуюся
func (t *Task) MarkAsInProgress() {
	t.Status = StatusInProgress
	t.UpdatedAt = time.Now()
}

// MarkAsCompleted помечает задачу как завершенную
func (t *Task) MarkAsCompleted() {
	t.Status = StatusCompleted
	t.UpdatedAt = time.Now()
}

// MarkAsCancelled помечает задачу как отмененную
func (t *Task) MarkAsCancelled() {
	t.Status = StatusCancelled
	t.UpdatedAt = time.Now()
}

// IsCompleted проверяет завершена ли задача
func (t *Task) IsCompleted() bool {
	return t.Status == StatusCompleted
}

// IsOverdue проверяет просрочена ли задача
func (t *Task) IsOverdue() bool {
	if t.DueDate == nil {
		return false
	}
	return time.Now().After(*t.DueDate) && !t.IsCompleted()
}

// BelongsToUser проверяет принадлежит ли задача пользователю
func (t *Task) BelongsToUser(userID uint) bool {
	return t.UserID == userID
}

// Константы для статусов задач
const (
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusCancelled  = "cancelled"
)

// Константы для приоритетов
const (
	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
)

// isValidTaskStatus проверяет корректность статуса задачи
func isValidTaskStatus(status string) bool {
	return status == StatusPending ||
		status == StatusInProgress ||
		status == StatusCompleted ||
		status == StatusCancelled
}

// isValidPriority проверяет корректность приоритета
func isValidPriority(priority string) bool {
	return priority == PriorityLow ||
		priority == PriorityMedium ||
		priority == PriorityHigh
}

