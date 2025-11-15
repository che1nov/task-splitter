package usecases

import (
	"context"
	"task-splitter/internal/dto"
	"task-splitter/pkg/logger"
)

// GetUserUseCase Use Case для получения пользователя
type GetUserUseCase struct {
	userRepo UserRepository
	log      logger.Logger
}

// NewGetUserUseCase создает новый UseCase для получения пользователя
func NewGetUserUseCase(userRepo UserRepository, log logger.Logger) *GetUserUseCase {
	return &GetUserUseCase{
		userRepo: userRepo,
		log:      log,
	}
}

// Execute получает пользователя по ID
func (uc *GetUserUseCase) Execute(ctx context.Context, userID uint) (dto.UserOutput, error) {
	uc.log.DebugContext(ctx, "Получение пользователя", "user_id", userID)

	// Получаем пользователя из БД
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		uc.log.ErrorContext(ctx, "Ошибка получения пользователя", "user_id", userID, "error", err)
		return dto.UserOutput{}, err
	}

	uc.log.DebugContext(ctx, "Пользователь получен", "user_id", userID, "username", user.Username)

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

