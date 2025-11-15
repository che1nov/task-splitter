package http

import (
	"net/http"
	"task-splitter/internal/dto"
	"task-splitter/internal/usecases"
	"task-splitter/pkg/logger"

	"github.com/gin-gonic/gin"
)

// AuthHandler контроллер для аутентификации
type AuthHandler struct {
	loginUC    *usecases.LoginUseCase
	registerUC *usecases.RegisterUseCase
	log        logger.Logger
}

// NewAuthHandler создает новый контроллер аутентификации
func NewAuthHandler(
	loginUC *usecases.LoginUseCase,
	registerUC *usecases.RegisterUseCase,
	log logger.Logger,
) *AuthHandler {
	return &AuthHandler{
		loginUC:    loginUC,
		registerUC: registerUC,
		log:        log,
	}
}

// Login выполняет вход пользователя
// @Summary User login
// @Description Authenticate user and return token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginInput true "Login credentials"
// @Success 200 {object} dto.LoginOutput
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.loginUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, output)
}

// Register регистрирует нового пользователя
// @Summary User registration
// @Description Register new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterInput true "Registration data"
// @Success 201 {object} dto.RegisterOutput
// @Failure 400 {object} gin.H
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input dto.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.registerUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, output)
}

