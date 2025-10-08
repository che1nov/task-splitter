package models

import (
	"time"

	"gorm.io/gorm"
)

// User представляет пользователя системы
type User struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	KeycloakID string         `json:"keycloak_id" gorm:"uniqueIndex;not null"`
	Email      string         `json:"email" gorm:"uniqueIndex;not null"`
	Username   string         `json:"username" gorm:"uniqueIndex;not null"`
	FirstName  string         `json:"first_name"`
	LastName   string         `json:"last_name"`
	Role       string         `json:"role" gorm:"default:'free'"` // free, premium
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Связи
	Tasks []Task `json:"tasks,omitempty" gorm:"foreignKey:UserID"`
}

// Task представляет основную задачу
type Task struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	Status      string         `json:"status" gorm:"default:'pending'"`  // pending, in_progress, completed, cancelled
	Priority    string         `json:"priority" gorm:"default:'medium'"` // low, medium, high
	DueDate     *time.Time     `json:"due_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Связи
	User      User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Subtasks  []Subtask  `json:"subtasks,omitempty" gorm:"foreignKey:TaskID"`
	Templates []Template `json:"templates,omitempty" gorm:"many2many:task_templates;"`
}

// Subtask представляет подзадачу
type Subtask struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	TaskID        uint           `json:"task_id" gorm:"not null"`
	Title         string         `json:"title" gorm:"not null"`
	Description   string         `json:"description"`
	Status        string         `json:"status" gorm:"default:'pending'"`  // pending, in_progress, completed
	Priority      string         `json:"priority" gorm:"default:'medium'"` // low, medium, high
	Order         int            `json:"order" gorm:"default:0"`
	EstimatedTime *int           `json:"estimated_time"` // в минутах
	ActualTime    *int           `json:"actual_time"`    // в минутах
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Связи
	Task Task `json:"task,omitempty" gorm:"foreignKey:TaskID"`
}

// Template представляет шаблон для разбивки задач
type Template struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Category    string         `json:"category"` // work, personal, project, etc.
	Prompt      string         `json:"prompt" gorm:"type:text"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Связи
	Tasks []Task `json:"tasks,omitempty" gorm:"many2many:task_templates;"`
}

// TaskSplitRequest представляет запрос на разбивку задачи
type TaskSplitRequest struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	TaskID     uint      `json:"task_id" gorm:"not null"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	Text       string    `json:"text" gorm:"type:text"`
	TemplateID *uint     `json:"template_id"`
	Status     string    `json:"status" gorm:"default:'pending'"` // pending, processing, completed, failed
	Result     string    `json:"result" gorm:"type:text"`
	Error      string    `json:"error"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Связи
	Task     Task      `json:"task,omitempty" gorm:"foreignKey:TaskID"`
	User     User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Template *Template `json:"template,omitempty" gorm:"foreignKey:TemplateID"`
}

// CreateUserRequest представляет запрос на создание пользователя
type CreateUserRequest struct {
	KeycloakID string `json:"keycloak_id" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Username   string `json:"username" binding:"required"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Role       string `json:"role"`
}

// UpdateUserRequest представляет запрос на обновление пользователя
type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}

// CreateTaskRequest представляет запрос на создание задачи
type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
}

// UpdateTaskRequest представляет запрос на обновление задачи
type UpdateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
}

// SplitTaskRequest представляет запрос на разбивку задачи
type SplitTaskRequest struct {
	TaskID     uint   `json:"task_id" binding:"required"`
	Text       string `json:"text" binding:"required"`
	TemplateID *uint  `json:"template_id"`
}

// SplitTaskResponse представляет ответ на разбивку задачи
type SplitTaskResponse struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// SplitStatusResponse представляет статус разбивки задачи
type SplitStatusResponse struct {
	RequestID string    `json:"request_id"`
	Status    string    `json:"status"`
	Result    string    `json:"result,omitempty"`
	Error     string    `json:"error,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
