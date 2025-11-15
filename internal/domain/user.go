package domain

import (
	"regexp"
	"strings"
	"time"
)

// User представляет пользователя системы
type User struct {
	ID         uint
	KeycloakID string
	Email      string
	Username   string
	FirstName  string
	LastName   string
	Role       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewUser создает нового пользователя с валидацией
func NewUser(keycloakID, email, username, firstName, lastName, role string) (User, error) {
	// Валидация KeycloakID
	if strings.TrimSpace(keycloakID) == "" {
		return User{}, ErrEmptyKeycloakID
	}

	// Валидация Email
	if !isValidEmail(email) {
		return User{}, ErrInvalidEmail
	}

	// Валидация Username
	if !isValidUsername(username) {
		return User{}, ErrInvalidUsername
	}

	// Валидация Role
	if role == "" {
		role = RoleFree // Роль по умолчанию
	}
	if !isValidRole(role) {
		return User{}, ErrInvalidRole
	}

	return User{
		KeycloakID: strings.TrimSpace(keycloakID),
		Email:      strings.ToLower(strings.TrimSpace(email)),
		Username:   strings.TrimSpace(username),
		FirstName:  strings.TrimSpace(firstName),
		LastName:   strings.TrimSpace(lastName),
		Role:       role,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

// UpdateProfile обновляет профиль пользователя
func (u *User) UpdateProfile(firstName, lastName string) {
	u.FirstName = strings.TrimSpace(firstName)
	u.LastName = strings.TrimSpace(lastName)
	u.UpdatedAt = time.Now()
}

// ChangeRole изменяет роль пользователя
func (u *User) ChangeRole(role string) error {
	if !isValidRole(role) {
		return ErrInvalidRole
	}
	u.Role = role
	u.UpdatedAt = time.Now()
	return nil
}

// IsPremium проверяет является ли пользователь премиум
func (u *User) IsPremium() bool {
	return u.Role == RolePremium
}

// IsAdmin проверяет является ли пользователь администратором
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// FullName возвращает полное имя пользователя
func (u *User) FullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.Username
	}
	return strings.TrimSpace(u.FirstName + " " + u.LastName)
}

// Константы для ролей пользователей
const (
	RoleFree    = "free"
	RolePremium = "premium"
	RoleAdmin   = "admin"
)

// isValidEmail проверяет корректность email
func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}
	
	// Простая проверка формата email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// isValidUsername проверяет корректность username
func isValidUsername(username string) bool {
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 50 {
		return false
	}
	
	// Username может содержать буквы, цифры, подчеркивания и дефисы
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
	return usernameRegex.MatchString(username)
}

// isValidRole проверяет корректность роли
func isValidRole(role string) bool {
	return role == RoleFree || role == RolePremium || role == RoleAdmin
}

