package postgresql

import (
	"time"

	"gorm.io/gorm"
)

// UserModel модель пользователя для GORM
type UserModel struct {
	ID         uint           `gorm:"primaryKey"`
	KeycloakID string         `gorm:"uniqueIndex;not null"`
	Email      string         `gorm:"uniqueIndex;not null"`
	Username   string         `gorm:"uniqueIndex;not null"`
	FirstName  string         `gorm:""`
	LastName   string         `gorm:""`
	Role       string         `gorm:"default:'free'"`
	CreatedAt  time.Time      `gorm:""`
	UpdatedAt  time.Time      `gorm:""`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// TableName указывает имя таблицы для UserModel
func (UserModel) TableName() string {
	return "users"
}

// TaskModel модель задачи для GORM
type TaskModel struct {
	ID          uint           `gorm:"primaryKey"`
	UserID      uint           `gorm:"not null;index"`
	Title       string         `gorm:"not null"`
	Description string         `gorm:""`
	Status      string         `gorm:"default:'pending'"`
	Priority    string         `gorm:"default:'medium'"`
	DueDate     *time.Time     `gorm:""`
	CreatedAt   time.Time      `gorm:""`
	UpdatedAt   time.Time      `gorm:""`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName указывает имя таблицы для TaskModel
func (TaskModel) TableName() string {
	return "tasks"
}

// SubtaskModel модель подзадачи для GORM
type SubtaskModel struct {
	ID            uint           `gorm:"primaryKey"`
	TaskID        uint           `gorm:"not null;index"`
	Title         string         `gorm:"not null"`
	Description   string         `gorm:""`
	Status        string         `gorm:"default:'pending'"`
	Priority      string         `gorm:"default:'medium'"`
	Order         int            `gorm:"default:0"`
	EstimatedTime *int           `gorm:""`
	ActualTime    *int           `gorm:""`
	CreatedAt     time.Time      `gorm:""`
	UpdatedAt     time.Time      `gorm:""`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// TableName указывает имя таблицы для SubtaskModel
func (SubtaskModel) TableName() string {
	return "subtasks"
}

// TemplateModel модель шаблона для GORM
type TemplateModel struct {
	ID          uint           `gorm:"primaryKey"`
	Name        string         `gorm:"not null"`
	Description string         `gorm:""`
	Category    string         `gorm:""`
	Prompt      string         `gorm:"type:text"`
	IsActive    bool           `gorm:"default:true"`
	CreatedAt   time.Time      `gorm:""`
	UpdatedAt   time.Time      `gorm:""`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName указывает имя таблицы для TemplateModel
func (TemplateModel) TableName() string {
	return "templates"
}

// TaskSplitRequestModel модель запроса на разбивку для GORM
type TaskSplitRequestModel struct {
	ID         string    `gorm:"primaryKey"`
	TaskID     uint      `gorm:"not null;index"`
	UserID     uint      `gorm:"not null;index"`
	Text       string    `gorm:"type:text"`
	TemplateID *uint     `gorm:""`
	Status     string    `gorm:"default:'pending'"`
	Result     string    `gorm:"type:text"`
	Error      string    `gorm:""`
	CreatedAt  time.Time `gorm:""`
	UpdatedAt  time.Time `gorm:""`
}

// TableName указывает имя таблицы для TaskSplitRequestModel
func (TaskSplitRequestModel) TableName() string {
	return "task_split_requests"
}

