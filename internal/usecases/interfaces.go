package usecases

import (
	"context"
	"task-splitter/internal/domain"
)

// UserRepository интерфейс адаптера для работы с пользователями
type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) error
	GetUserByID(ctx context.Context, id uint) (domain.User, error)
	GetUserByKeycloakID(ctx context.Context, keycloakID string) (domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) error
	DeleteUser(ctx context.Context, id uint) error
	ListUsers(ctx context.Context, limit, offset int) ([]domain.User, error)
}

// TaskRepository интерфейс адаптера для работы с задачами
type TaskRepository interface {
	CreateTask(ctx context.Context, task domain.Task) error
	GetTaskByID(ctx context.Context, id uint) (domain.Task, error)
	GetTasksByUserID(ctx context.Context, userID uint, limit, offset int) ([]domain.Task, error)
	GetTaskWithSubtasks(ctx context.Context, id uint) (domain.Task, []domain.Subtask, error)
	UpdateTask(ctx context.Context, task domain.Task) error
	DeleteTask(ctx context.Context, id uint) error
}

// SubtaskRepository интерфейс адаптера для работы с подзадачами
type SubtaskRepository interface {
	CreateSubtask(ctx context.Context, subtask domain.Subtask) error
	GetSubtaskByID(ctx context.Context, id uint) (domain.Subtask, error)
	GetSubtasksByTaskID(ctx context.Context, taskID uint) ([]domain.Subtask, error)
	UpdateSubtask(ctx context.Context, subtask domain.Subtask) error
	DeleteSubtask(ctx context.Context, id uint) error
	DeleteSubtasksByTaskID(ctx context.Context, taskID uint) error
}

// TemplateRepository интерфейс адаптера для работы с шаблонами
type TemplateRepository interface {
	CreateTemplate(ctx context.Context, template domain.Template) error
	GetTemplateByID(ctx context.Context, id uint) (domain.Template, error)
	GetTemplatesByCategory(ctx context.Context, category string) ([]domain.Template, error)
	GetActiveTemplates(ctx context.Context) ([]domain.Template, error)
	UpdateTemplate(ctx context.Context, template domain.Template) error
	DeleteTemplate(ctx context.Context, id uint) error
	ListTemplates(ctx context.Context, limit, offset int) ([]domain.Template, error)
}

// TaskSplitRequestRepository интерфейс адаптера для работы с запросами разбивки
type TaskSplitRequestRepository interface {
	CreateSplitRequest(ctx context.Context, request domain.TaskSplitRequest) error
	GetSplitRequestByID(ctx context.Context, id string) (domain.TaskSplitRequest, error)
	GetSplitRequestsByUserID(ctx context.Context, userID uint, limit, offset int) ([]domain.TaskSplitRequest, error)
	GetSplitRequestsByStatus(ctx context.Context, status string) ([]domain.TaskSplitRequest, error)
	UpdateSplitRequest(ctx context.Context, request domain.TaskSplitRequest) error
	DeleteSplitRequest(ctx context.Context, id string) error
}

// CacheAdapter интерфейс адаптера для кэширования
type CacheAdapter interface {
	Set(ctx context.Context, key string, value interface{}, ttlSeconds int) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

// MessagingAdapter интерфейс адаптера для обмена сообщениями
type MessagingAdapter interface {
	PublishSplitTaskRequest(ctx context.Context, requestID string, taskID, userID uint, text string, templateID *uint) error
	Close() error
}

// KeycloakAdapter интерфейс адаптера для Keycloak
type KeycloakAdapter interface {
	ValidateToken(ctx context.Context, token string) (string, error) // Возвращает keycloakID
	GetUserInfo(ctx context.Context, keycloakID string) (KeycloakUserInfo, error)
}

// KeycloakUserInfo информация о пользователе из Keycloak
type KeycloakUserInfo struct {
	Sub               string
	Email             string
	PreferredUsername string
	GivenName         string
	FamilyName        string
	RealmRoles        []string
	Groups            []string
}

