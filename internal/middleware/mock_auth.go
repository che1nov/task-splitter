package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MockAuthMiddleware создает простой mock middleware для демо
func MockAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Проверяем формат Bearer token
		if authHeader != "Bearer mock_token_user123" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Сохраняем mock информацию о пользователе в контексте
		c.Set("user_id", uint(1))
		c.Set("user_sub", "user123")
		c.Set("user_email", "test@example.com")
		c.Set("user_roles", []string{"free"})

		c.Next()
	}
}
