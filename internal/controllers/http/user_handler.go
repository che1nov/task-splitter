package http

import (
	"net/http"
	"task-splitter/internal/dto"
	"task-splitter/internal/usecases"
	"task-splitter/pkg/logger"

	"github.com/gin-gonic/gin"
)

// UserHandler контроллер для работы с пользователями
type UserHandler struct {
	getUserUC    *usecases.GetUserUseCase
	updateUserUC *usecases.UpdateUserUseCase
	log          logger.Logger
}

// NewUserHandler создает новый контроллер пользователей
func NewUserHandler(
	getUserUC *usecases.GetUserUseCase,
	updateUserUC *usecases.UpdateUserUseCase,
	log logger.Logger,
) *UserHandler {
	return &UserHandler{
		getUserUC:    getUserUC,
		updateUserUC: updateUserUC,
		log:          log,
	}
}

// GetProfile получает профиль текущего пользователя
// @Summary Get user profile
// @Description Get current user profile information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserOutput
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Получаем ID пользователя из контекста (устанавливается middleware)
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Вызываем Use Case
	output, err := h.getUserUC.Execute(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, output)
}

// UpdateProfile обновляет профиль текущего пользователя
// @Summary Update user profile
// @Description Update current user profile information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateUserInput true "User update data"
// @Success 200 {object} dto.UserOutput
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// Получаем ID пользователя из контекста
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Парсим тело запроса
	var input dto.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Вызываем Use Case
	output, err := h.updateUserUC.Execute(c.Request.Context(), userID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}

