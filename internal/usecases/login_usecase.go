package usecases

import (
	"context"
	"task-splitter/internal/dto"
	"task-splitter/pkg/logger"
)

// LoginUseCase Use Case для входа пользователя
type LoginUseCase struct {
	userRepo UserRepository
	log      logger.Logger
}

// NewLoginUseCase создает новый UseCase для входа
func NewLoginUseCase(userRepo UserRepository, log logger.Logger) *LoginUseCase {
	return &LoginUseCase{
		userRepo: userRepo,
		log:      log,
	}
}

// Execute выполняет вход пользователя
func (uc *LoginUseCase) Execute(ctx context.Context, input dto.LoginInput) (dto.LoginOutput, error) {
	uc.log.InfoContext(ctx, "Попытка входа", "username", input.Username)

	// Получаем пользователя по username
	user, err := uc.userRepo.GetUserByUsername(ctx, input.Username)
	if err != nil {
		uc.log.WarnContext(ctx, "Пользователь не найден", "username", input.Username)
		return dto.LoginOutput{}, err
	}

	// В реальном приложении здесь должна быть проверка пароля через Keycloak
	// Пока что просто возвращаем mock токен

	uc.log.InfoContext(ctx, "Вход выполнен успешно", "user_id", user.ID, "username", user.Username)

	return dto.LoginOutput{
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
		Token: "mock_token_" + user.KeycloakID,
	}, nil
}

