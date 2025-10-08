package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"task-splitter/pkg/keycloak"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware создает middleware для аутентификации через Keycloak
func AuthMiddleware(keycloakClient *keycloak.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Проверяем формат Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Валидируем токен через Keycloak
		tokenInfo, err := keycloakClient.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		// Получаем информацию о пользователе
		userInfo, err := keycloakClient.GetUserInfo(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to get user info", "details": err.Error()})
			c.Abort()
			return
		}

		// Сохраняем информацию о пользователе в контексте
		c.Set("user_id", tokenInfo.Sub)
		c.Set("user_email", userInfo.Email)
		c.Set("user_username", userInfo.PreferredUsername)
		c.Set("user_roles", tokenInfo.RealmRoles)
		c.Set("user_groups", tokenInfo.Groups)

		c.Next()
	}
}

// RoleMiddleware создает middleware для проверки ролей
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoles, exists := c.Get("user_roles")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "User roles not found"})
			c.Abort()
			return
		}

		roles, ok := userRoles.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid user roles format"})
			c.Abort()
			return
		}

		// Проверяем, есть ли у пользователя хотя бы одна из требуемых ролей
		hasRequiredRole := false
		for _, requiredRole := range requiredRoles {
			for _, userRole := range roles {
				if userRole == requiredRole {
					hasRequiredRole = true
					break
				}
			}
			if hasRequiredRole {
				break
			}
		}

		if !hasRequiredRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error":          "Insufficient permissions",
				"required_roles": requiredRoles,
				"user_roles":     roles,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// PremiumMiddleware проверяет, что пользователь имеет premium роль
func PremiumMiddleware() gin.HandlerFunc {
	return RoleMiddleware("premium")
}

// CORS middleware для обработки CORS запросов
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Logger middleware для логирования запросов
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("02/Jan/2006:15:04:05 -0700"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// RateLimitMiddleware создает middleware для ограничения скорости запросов
func RateLimitMiddleware() gin.HandlerFunc {
	// Простая реализация rate limiting
	// В продакшене лучше использовать Redis или специализированные библиотеки
	return func(c *gin.Context) {
		// Здесь можно добавить логику rate limiting
		c.Next()
	}
}

// RequestIDMiddleware добавляет уникальный ID к каждому запросу
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// generateRequestID генерирует уникальный ID запроса
func generateRequestID() string {
	// Простая реализация генерации ID
	// В продакшене лучше использовать UUID
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
