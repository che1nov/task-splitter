package usecases

import (
	"context"
	"fmt"
	"task-splitter/internal/domain"
	"task-splitter/internal/dto"
	"task-splitter/pkg/logger"
	"time"
)

// RegisterUseCase Use Case для регистрации пользователя
type RegisterUseCase struct {
	userRepo UserRepository
	log      logger.Logger
}

// NewRegisterUseCase создает новый UseCase для регистрации
func NewRegisterUseCase(userRepo UserRepository, log logger.Logger) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo: userRepo,
		log:      log,
	}
}

// Execute регистрирует нового пользователя
func (uc *RegisterUseCase) Execute(ctx context.Context, input dto.RegisterInput) (dto.RegisterOutput, error) {
	uc.log.InfoContext(ctx, "Регистрация нового пользователя", "username", input.Username, "email", input.Email)

	// Генерируем KeycloakID для демо-режима
	keycloakID := "demo_" + input.Username + "_" + fmt.Sprintf("%d", time.Now().Unix())

	// Создаем доменную сущность с валидацией
	user, err := domain.NewUser(keycloakID, input.Email, input.Username, input.FirstName, input.LastName, domain.RoleFree)
	if err != nil {
		uc.log.ErrorContext(ctx, "Ошибка валидации при регистрации", "error", err)
		return dto.RegisterOutput{}, err
	}

	// Сохраняем в БД
	if err := uc.userRepo.CreateUser(ctx, user); err != nil {
		uc.log.ErrorContext(ctx, "Ошибка создания пользователя в БД", "error", err)
		return dto.RegisterOutput{}, err
	}

	uc.log.InfoContext(ctx, "Пользователь успешно зарегистрирован", "user_id", user.ID, "username", user.Username)

	return dto.RegisterOutput{
		Message: "Пользователь успешно зарегистрирован",
		User: dto.UserOutput{
			ID:         user.ID,
			KeycloakID: user.KeycloakID,
			Email:      user.Email,
			Username:   user.Username,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Role:       user.Role,
			CreatedAt:  user.CreatedAt,
			UpdatedAt:  user.UpdatedAt,
		},
	}, nil
}

