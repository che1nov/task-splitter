package dto

import "time"

// CreateTaskInput входные данные для создания задачи
type CreateTaskInput struct {
	UserID      uint       `json:"-"` // Берется из контекста
	Title       string     `json:"title" validate:"required,min=1,max=500"`
	Description string     `json:"description" validate:"max=5000"`
	Priority    string     `json:"priority" validate:"omitempty,oneof=low medium high"`
	DueDate     *time.Time `json:"due_date"`
}

// UpdateTaskInput входные данные для обновления задачи
type UpdateTaskInput struct {
	TaskID      uint       `json:"-"` // Из URL
	Title       string     `json:"title" validate:"omitempty,min=1,max=500"`
	Description string     `json:"description" validate:"max=5000"`
	Status      string     `json:"status" validate:"omitempty,oneof=pending in_progress completed cancelled"`
	Priority    string     `json:"priority" validate:"omitempty,oneof=low medium high"`
	DueDate     *time.Time `json:"due_date"`
}

// TaskOutput выходные данные задачи
type TaskOutput struct {
	ID          uint          `json:"id"`
	UserID      uint          `json:"user_id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      string        `json:"status"`
	Priority    string        `json:"priority"`
	DueDate     *time.Time    `json:"due_date"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Subtasks    []SubtaskOutput `json:"subtasks,omitempty"`
}

// GetTaskInput входные данные для получения задачи
type GetTaskInput struct {
	TaskID uint `json:"task_id" validate:"required"`
}

// GetTasksByUserInput входные данные для получения задач пользователя
type GetTasksByUserInput struct {
	UserID uint `json:"user_id" validate:"required"`
	Limit  int  `json:"limit" validate:"min=1,max=100"`
	Offset int  `json:"offset" validate:"min=0"`
}

// DeleteTaskInput входные данные для удаления задачи
type DeleteTaskInput struct {
	TaskID uint `json:"task_id" validate:"required"`
	UserID uint `json:"user_id" validate:"required"`
}

// SubtaskOutput выходные данные подзадачи
type SubtaskOutput struct {
	ID            uint      `json:"id"`
	TaskID        uint      `json:"task_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	Priority      string    `json:"priority"`
	Order         int       `json:"order"`
	EstimatedTime *int      `json:"estimated_time"`
	ActualTime    *int      `json:"actual_time"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// SplitTaskInput входные данные для разбивки задачи
type SplitTaskInput struct {
	UserID     uint   `json:"-"` // Из контекста
	TaskID     uint   `json:"task_id" validate:"required"`
	Text       string `json:"text" validate:"required,min=10,max=10000"`
	TemplateID *uint  `json:"template_id"`
}

// SplitTaskOutput выходные данные после отправки на разбивку
type SplitTaskOutput struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// GetSplitStatusInput входные данные для получения статуса разбивки
type GetSplitStatusInput struct {
	RequestID string `json:"request_id" validate:"required"`
}

// SplitStatusOutput выходные данные статуса разбивки
type SplitStatusOutput struct {
	RequestID string    `json:"request_id"`
	Status    string    `json:"status"`
	Result    string    `json:"result,omitempty"`
	Error     string    `json:"error,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

