package usecases

import (
	"context"
	"task-splitter/internal/domain"
	"task-splitter/internal/dto"
	"task-splitter/pkg/logger"
)

// CreateUserUseCase Use Case для создания пользователя
type CreateUserUseCase struct {
	userRepo UserRepository
	log      logger.Logger
}

// NewCreateUserUseCase создает новый UseCase для создания пользователя
func NewCreateUserUseCase(userRepo UserRepository, log logger.Logger) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo: userRepo,
		log:      log,
	}
}

// Execute создает нового пользователя
func (uc *CreateUserUseCase) Execute(ctx context.Context, input dto.CreateUserInput) (dto.UserOutput, error) {
	uc.log.InfoContext(ctx, "Создание пользователя", "email", input.Email, "username", input.Username)

	// Создаем доменную сущность с валидацией
	user, err := domain.NewUser(input.KeycloakID, input.Email, input.Username, input.FirstName, input.LastName, input.Role)
	if err != nil {
		uc.log.ErrorContext(ctx, "Ошибка валидации пользователя", "error", err)
		return dto.UserOutput{}, err
	}

	// Сохраняем в БД
	if err := uc.userRepo.CreateUser(ctx, user); err != nil {
		uc.log.ErrorContext(ctx, "Ошибка сохранения пользователя в БД", "error", err)
		return dto.UserOutput{}, err
	}

	uc.log.InfoContext(ctx, "Пользователь успешно создан", "id", user.ID, "username", user.Username)

	return dto.UserOutput{
		ID:         user.ID,
		KeycloakID: user.KeycloakID,
		Email:      user.Email,
		Username:   user.Username,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Role:       user.Role,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}, nil
}

