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
	"gorm.io/gorm"
	redisLib "github.com/redis/go-redis/v9"
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

	// Настройка Gin
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// Health check - должен работать независимо от других сервисов
	healthHandler := handlers.NewHealthHandler()
	r.GET("/api/v1/health", healthHandler.HealthCheck)

	// Попытка инициализации базы данных
	var db *gorm.DB
	var err error
	
	log.Println("Attempting to connect to database...")
	db, err = postgres.NewConnection(cfg.Database)
	if err != nil {
		log.Printf("⚠️  Failed to connect to database: %v", err)
		log.Println("⚠️  Server will start in limited mode (health check only)")
		db = nil
	} else {
		log.Println("✅ Database connected successfully")
	}

	// Попытка инициализации Redis
	var redisClient *redisLib.Client
	log.Println("Attempting to connect to Redis...")
	redisClient = redis.NewClient(cfg.Redis)
	if redisClient == nil {
		log.Println("⚠️  Redis connection failed, server will continue without Redis")
	}

	// Попытка инициализации Keycloak
	var keycloakClient *keycloak.Client
	log.Println("Attempting to connect to Keycloak...")
	keycloakClient = keycloak.NewClient(cfg.Keycloak)
	if keycloakClient == nil {
		log.Println("⚠️  Keycloak connection failed, server will continue without Keycloak")
	}

	// Инициализация handlers только если база данных доступна
	if db != nil {
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
	} else {
		log.Println("⚠️  Running in limited mode - only health check endpoint available")
	}

	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("🚀 Server starting on %s", serverAddr)
	log.Printf("🌐 Health check available at: http://%s/api/v1/health", serverAddr)
	log.Fatal(r.Run(serverAddr))
}
