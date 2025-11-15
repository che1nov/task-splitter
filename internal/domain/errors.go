package domain

import "errors"

// Доменные ошибки для User
var (
	ErrUserNotFound      = errors.New("пользователь не найден")
	ErrUserAlreadyExists = errors.New("пользователь уже существует")
	ErrInvalidEmail      = errors.New("некорректный email")
	ErrInvalidUsername   = errors.New("некорректный username")
	ErrInvalidRole       = errors.New("некорректная роль")
	ErrEmptyKeycloakID   = errors.New("keycloak ID не может быть пустым")
)

// Доменные ошибки для Task
var (
	ErrTaskNotFound      = errors.New("задача не найдена")
	ErrInvalidTitle      = errors.New("название задачи не может быть пустым")
	ErrInvalidPriority   = errors.New("некорректный приоритет")
	ErrInvalidStatus     = errors.New("некорректный статус задачи")
	ErrAccessDenied      = errors.New("доступ запрещен")
)

// Доменные ошибки для Subtask
var (
	ErrSubtaskNotFound        = errors.New("подзадача не найдена")
	ErrInvalidSubtaskTitle    = errors.New("название подзадачи не может быть пустым")
	ErrInvalidSubtaskStatus   = errors.New("некорректный статус подзадачи")
	ErrInvalidSubtaskPriority = errors.New("некорректный приоритет подзадачи")
)

// Доменные ошибки для Template
var (
	ErrTemplateNotFound   = errors.New("шаблон не найден")
	ErrInvalidTemplate    = errors.New("некорректный шаблон")
	ErrEmptyTemplateName  = errors.New("название шаблона не может быть пустым")
	ErrEmptyTemplatePrompt = errors.New("промпт шаблона не может быть пустым")
)

// Доменные ошибки для TaskSplitRequest
var (
	ErrSplitRequestNotFound = errors.New("запрос на разбивку не найден")
	ErrEmptySplitText       = errors.New("текст для разбивки не может быть пустым")
	ErrInvalidSplitStatus   = errors.New("некорректный статус разбивки")
)

