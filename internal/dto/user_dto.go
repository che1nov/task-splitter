package dto

import "time"

// CreateUserInput входные данные для создания пользователя
type CreateUserInput struct {
	KeycloakID string `json:"keycloak_id" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Username   string `json:"username" validate:"required,min=3,max=50"`
	FirstName  string `json:"first_name" validate:"max=100"`
	LastName   string `json:"last_name" validate:"max=100"`
	Role       string `json:"role" validate:"oneof=free premium admin"`
}

// UpdateUserInput входные данные для обновления пользователя
type UpdateUserInput struct {
	FirstName string `json:"first_name" validate:"max=100"`
	LastName  string `json:"last_name" validate:"max=100"`
	Role      string `json:"role" validate:"oneof=free premium admin"`
}

// UserOutput выходные данные пользователя
type UserOutput struct {
	ID         uint      `json:"id"`
	KeycloakID string    `json:"keycloak_id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// GetUserByIDInput входные данные для получения пользователя по ID
type GetUserByIDInput struct {
	UserID uint `json:"user_id" validate:"required"`
}

// GetUserByKeycloakIDInput входные данные для получения пользователя по Keycloak ID
type GetUserByKeycloakIDInput struct {
	KeycloakID string `json:"keycloak_id" validate:"required"`
}

// GetUserByUsernameInput входные данные для получения пользователя по username
type GetUserByUsernameInput struct {
	Username string `json:"username" validate:"required"`
}

