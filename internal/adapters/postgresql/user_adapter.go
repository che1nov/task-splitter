package postgresql

import (
	"context"
	"task-splitter/internal/domain"
	"task-splitter/pkg/logger"

	"gorm.io/gorm"
)

// UserAdapter адаптер для работы с пользователями в PostgreSQL
type UserAdapter struct {
	db  *gorm.DB
	log logger.Logger
}

// NewUserAdapter создает новый адаптер пользователей
func NewUserAdapter(db *gorm.DB, log logger.Logger) *UserAdapter {
	return &UserAdapter{
		db:  db,
		log: log,
	}
}

// CreateUser создает пользователя в БД
func (a *UserAdapter) CreateUser(ctx context.Context, user domain.User) error {
	model := UserModel{
		KeycloakID: user.KeycloakID,
		Email:      user.Email,
		Username:   user.Username,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Role:       user.Role,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}

	if err := a.db.WithContext(ctx).Create(&model).Error; err != nil {
		a.log.ErrorContext(ctx, "Ошибка создания пользователя в БД", "error", err)
		return err
	}

	return nil
}

// GetUserByID получает пользователя по ID
func (a *UserAdapter) GetUserByID(ctx context.Context, id uint) (domain.User, error) {
	var model UserModel
	if err := a.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.User{}, domain.ErrUserNotFound
		}
		a.log.ErrorContext(ctx, "Ошибка получения пользователя по ID", "id", id, "error", err)
		return domain.User{}, err
	}

	return toDomainUser(model), nil
}

// GetUserByKeycloakID получает пользователя по Keycloak ID
func (a *UserAdapter) GetUserByKeycloakID(ctx context.Context, keycloakID string) (domain.User, error) {
	var model UserModel
	if err := a.db.WithContext(ctx).Where("keycloak_id = ?", keycloakID).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.User{}, domain.ErrUserNotFound
		}
		a.log.ErrorContext(ctx, "Ошибка получения пользователя по Keycloak ID", "keycloak_id", keycloakID, "error", err)
		return domain.User{}, err
	}

	return toDomainUser(model), nil
}

// GetUserByUsername получает пользователя по username
func (a *UserAdapter) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	var model UserModel
	if err := a.db.WithContext(ctx).Where("username = ?", username).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.User{}, domain.ErrUserNotFound
		}
		a.log.ErrorContext(ctx, "Ошибка получения пользователя по username", "username", username, "error", err)
		return domain.User{}, err
	}

	return toDomainUser(model), nil
}

// GetUserByEmail получает пользователя по email
func (a *UserAdapter) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var model UserModel
	if err := a.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.User{}, domain.ErrUserNotFound
		}
		a.log.ErrorContext(ctx, "Ошибка получения пользователя по email", "email", email, "error", err)
		return domain.User{}, err
	}

	return toDomainUser(model), nil
}

// UpdateUser обновляет пользователя
func (a *UserAdapter) UpdateUser(ctx context.Context, user domain.User) error {
	model := UserModel{
		ID:         user.ID,
		KeycloakID: user.KeycloakID,
		Email:      user.Email,
		Username:   user.Username,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Role:       user.Role,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}

	if err := a.db.WithContext(ctx).Save(&model).Error; err != nil {
		a.log.ErrorContext(ctx, "Ошибка обновления пользователя", "id", user.ID, "error", err)
		return err
	}

	return nil
}

// DeleteUser удаляет пользователя
func (a *UserAdapter) DeleteUser(ctx context.Context, id uint) error {
	if err := a.db.WithContext(ctx).Delete(&UserModel{}, id).Error; err != nil {
		a.log.ErrorContext(ctx, "Ошибка удаления пользователя", "id", id, "error", err)
		return err
	}

	return nil
}

// ListUsers получает список пользователей с пагинацией
func (a *UserAdapter) ListUsers(ctx context.Context, limit, offset int) ([]domain.User, error) {
	var models []UserModel
	if err := a.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		a.log.ErrorContext(ctx, "Ошибка получения списка пользователей", "error", err)
		return nil, err
	}

	users := make([]domain.User, len(models))
	for i, model := range models {
		users[i] = toDomainUser(model)
	}

	return users, nil
}

// toDomainUser конвертирует модель БД в доменную модель
func toDomainUser(model UserModel) domain.User {
	return domain.User{
		ID:         model.ID,
		KeycloakID: model.KeycloakID,
		Email:      model.Email,
		Username:   model.Username,
		FirstName:  model.FirstName,
		LastName:   model.LastName,
		Role:       model.Role,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}
}

