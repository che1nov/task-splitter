package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler контроллер для health check
type HealthHandler struct{}

// NewHealthHandler создает новый контроллер health check
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck проверяет состояние приложения
// @Summary Health check
// @Description Check application health status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"message":   "TaskSplitter API is running",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

