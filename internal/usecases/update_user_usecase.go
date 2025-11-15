package usecases

import (
	"context"
	"task-splitter/internal/dto"
	"task-splitter/pkg/logger"
)

// UpdateUserUseCase Use Case для обновления пользователя
type UpdateUserUseCase struct {
	userRepo UserRepository
	log      logger.Logger
}

// NewUpdateUserUseCase создает новый UseCase для обновления пользователя
func NewUpdateUserUseCase(userRepo UserRepository, log logger.Logger) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userRepo: userRepo,
		log:      log,
	}
}

// Execute обновляет профиль пользователя
func (uc *UpdateUserUseCase) Execute(ctx context.Context, userID uint, input dto.UpdateUserInput) (dto.UserOutput, error) {
	uc.log.InfoContext(ctx, "Обновление профиля пользователя", "user_id", userID)

	// Получаем пользователя из БД
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		uc.log.ErrorContext(ctx, "Пользователь не найден", "user_id", userID, "error", err)
		return dto.UserOutput{}, err
	}

	// Обновляем профиль
	user.UpdateProfile(input.FirstName, input.LastName)

	// Обновляем роль если указана
	if input.Role != "" {
		if err := user.ChangeRole(input.Role); err != nil {
			uc.log.ErrorContext(ctx, "Ошибка изменения роли", "user_id", userID, "error", err)
			return dto.UserOutput{}, err
		}
	}

	// Сохраняем изменения
	if err := uc.userRepo.UpdateUser(ctx, user); err != nil {
		uc.log.ErrorContext(ctx, "Ошибка сохранения изменений", "user_id", userID, "error", err)
		return dto.UserOutput{}, err
	}

	uc.log.InfoContext(ctx, "Профиль пользователя обновлен", "user_id", userID)

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

