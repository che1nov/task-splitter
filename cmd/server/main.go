package main

import (
	"log"
	"task-splitter/config"
	"task-splitter/internal/handlers"
	"task-splitter/internal/middleware"
	"task-splitter/internal/repository"
	"task-splitter/internal/service"
	"task-splitter/pkg/keycloak"
	"task-splitter/pkg/postgres"
	"task-splitter/pkg/redis"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title TaskSplitter API
// @version 1.0
// @description API для разбивки задач на подзадачи
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Инициализация базы данных
	db, err := postgres.NewConnection(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Инициализация Redis
	redisClient := redis.NewClient(cfg.Redis)

	// Инициализация Keycloak
	keycloakClient := keycloak.NewClient(cfg.Keycloak)

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	templateRepo := repository.NewTemplateRepository(db)
	splitRequestRepo := repository.NewTaskSplitRequestRepository(db)

	// Инициализация сервисов
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo, templateRepo, splitRequestRepo, redisClient)
	authService := service.NewAuthService(keycloakClient, userService)

	// Инициализация handlers
	userHandler := handlers.NewUserHandler(userService)
	taskHandler := handlers.NewTaskHandler(taskService)
	authHandler := handlers.NewAuthHandler(authService)

	// Настройка Gin
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// Swagger документация
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Публичные маршруты
	public := r.Group("/api/v1")
	{
		public.POST("/auth/login", authHandler.Login)
		public.POST("/auth/register", authHandler.Register)
	}

	// Защищенные маршруты
	protected := r.Group("/api/v1")
	protected.Use(middleware.MockAuthMiddleware())
	{
		// Пользователи
		protected.GET("/users/profile", userHandler.GetProfile)
		protected.PUT("/users/profile", userHandler.UpdateProfile)

		// Задачи
		protected.GET("/tasks", taskHandler.GetTasks)
		protected.POST("/tasks", taskHandler.CreateTask)
		protected.GET("/tasks/:id", taskHandler.GetTask)
		protected.PUT("/tasks/:id", taskHandler.UpdateTask)
		protected.DELETE("/tasks/:id", taskHandler.DeleteTask)

		// Разбивка задач
		protected.POST("/split", taskHandler.SplitTask)
		protected.GET("/split/:id/status", taskHandler.GetSplitStatus)
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(r.Run(":" + cfg.Server.Port))
}
