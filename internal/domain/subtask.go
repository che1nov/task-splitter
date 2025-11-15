package domain

import (
	"strings"
	"time"
)

// Subtask представляет подзадачу
type Subtask struct {
	ID            uint
	TaskID        uint
	Title         string
	Description   string
	Status        string
	Priority      string
	Order         int
	EstimatedTime *int // в минутах
	ActualTime    *int // в минутах
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewSubtask создает новую подзадачу с валидацией
func NewSubtask(taskID uint, title, description, priority string, order int, estimatedTime *int) (Subtask, error) {
	// Валидация заголовка
	title = strings.TrimSpace(title)
	if title == "" {
		return Subtask{}, ErrInvalidSubtaskTitle
	}

	// Валидация приоритета
	if priority == "" {
		priority = PriorityMedium // Приоритет по умолчанию
	}
	if !isValidPriority(priority) {
		return Subtask{}, ErrInvalidSubtaskPriority
	}

	return Subtask{
		TaskID:        taskID,
		Title:         title,
		Description:   strings.TrimSpace(description),
		Status:        StatusPending, // Статус по умолчанию
		Priority:      priority,
		Order:         order,
		EstimatedTime: estimatedTime,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

// Update обновляет подзадачу
func (s *Subtask) Update(title, description, status, priority string) error {
	// Валидация и обновление заголовка
	if title != "" {
		title = strings.TrimSpace(title)
		if title == "" {
			return ErrInvalidSubtaskTitle
		}
		s.Title = title
	}

	// Обновление описания
	if description != "" {
		s.Description = strings.TrimSpace(description)
	}

	// Валидация и обновление статуса
	if status != "" {
		if !isValidSubtaskStatus(status) {
			return ErrInvalidSubtaskStatus
		}
		s.Status = status
	}

	// Валидация и обновление приоритета
	if priority != "" {
		if !isValidPriority(priority) {
			return ErrInvalidSubtaskPriority
		}
		s.Priority = priority
	}

	s.UpdatedAt = time.Now()
	return nil
}

// MarkAsInProgress помечает подзадачу как выполняющуюся
func (s *Subtask) MarkAsInProgress() {
	s.Status = StatusInProgress
	s.UpdatedAt = time.Now()
}

// MarkAsCompleted помечает подзадачу как завершенную
func (s *Subtask) MarkAsCompleted() {
	s.Status = StatusCompleted
	s.UpdatedAt = time.Now()
}

// SetActualTime устанавливает фактическое время выполнения
func (s *Subtask) SetActualTime(minutes int) {
	s.ActualTime = &minutes
	s.UpdatedAt = time.Now()
}

// IsCompleted проверяет завершена ли подзадача
func (s *Subtask) IsCompleted() bool {
	return s.Status == StatusCompleted
}

// IsDelayed проверяет есть ли задержка выполнения
func (s *Subtask) IsDelayed() bool {
	if s.EstimatedTime == nil || s.ActualTime == nil {
		return false
	}
	return *s.ActualTime > *s.EstimatedTime
}

// Константы для статусов подзадач
const (
	SubtaskStatusPending    = "pending"
	SubtaskStatusInProgress = "in_progress"
	SubtaskStatusCompleted  = "completed"
)

// isValidSubtaskStatus проверяет корректность статуса подзадачи
func isValidSubtaskStatus(status string) bool {
	return status == SubtaskStatusPending ||
		status == SubtaskStatusInProgress ||
		status == SubtaskStatusCompleted
}

