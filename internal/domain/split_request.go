package domain

import (
	"strings"
	"time"
)

// TaskSplitRequest представляет запрос на разбивку задачи
type TaskSplitRequest struct {
	ID         string
	TaskID     uint
	UserID     uint
	Text       string
	TemplateID *uint
	Status     string
	Result     string
	Error      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewTaskSplitRequest создает новый запрос на разбивку с валидацией
func NewTaskSplitRequest(id string, taskID, userID uint, text string, templateID *uint) (TaskSplitRequest, error) {
	// Валидация текста
	text = strings.TrimSpace(text)
	if text == "" {
		return TaskSplitRequest{}, ErrEmptySplitText
	}

	return TaskSplitRequest{
		ID:         id,
		TaskID:     taskID,
		UserID:     userID,
		Text:       text,
		TemplateID: templateID,
		Status:     SplitStatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

// MarkAsProcessing помечает запрос как обрабатываемый
func (r *TaskSplitRequest) MarkAsProcessing() {
	r.Status = SplitStatusProcessing
	r.UpdatedAt = time.Now()
}

// MarkAsCompleted помечает запрос как завершенный
func (r *TaskSplitRequest) MarkAsCompleted(result string) {
	r.Status = SplitStatusCompleted
	r.Result = result
	r.Error = ""
	r.UpdatedAt = time.Now()
}

// MarkAsFailed помечает запрос как провалившийся
func (r *TaskSplitRequest) MarkAsFailed(errorMsg string) {
	r.Status = SplitStatusFailed
	r.Error = errorMsg
	r.UpdatedAt = time.Now()
}

// IsCompleted проверяет завершен ли запрос
func (r *TaskSplitRequest) IsCompleted() bool {
	return r.Status == SplitStatusCompleted
}

// IsFailed проверяет провалился ли запрос
func (r *TaskSplitRequest) IsFailed() bool {
	return r.Status == SplitStatusFailed
}

// IsProcessing проверяет обрабатывается ли запрос
func (r *TaskSplitRequest) IsProcessing() bool {
	return r.Status == SplitStatusProcessing
}

// Константы для статусов разбивки
const (
	SplitStatusPending    = "pending"
	SplitStatusProcessing = "processing"
	SplitStatusCompleted  = "completed"
	SplitStatusFailed     = "failed"
)

// isValidSplitStatus проверяет корректность статуса разбивки
func isValidSplitStatus(status string) bool {
	return status == SplitStatusPending ||
		status == SplitStatusProcessing ||
		status == SplitStatusCompleted ||
		status == SplitStatusFailed
}

