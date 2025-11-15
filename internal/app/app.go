package app

import (
	"task-splitter/config"
	"task-splitter/internal/adapters/postgresql"
	"task-splitter/internal/adapters/rabbitmq"
	redisAdapter "task-splitter/internal/adapters/redis"
	"task-splitter/internal/controllers/http"
	"task-splitter/internal/middleware"
	"task-splitter/internal/usecases"
	"task-splitter/pkg/logger"
	pgClient "task-splitter/pkg/postgres"
	redisClient "task-splitter/pkg/redis"

	"github.com/gin-gonic/gin"
	redisLib "github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// App представляет приложение
type App struct {
	router    *gin.Engine
	db        *gorm.DB
	redis     *redisLib.Client
	messaging *rabbitmq.MessagingAdapter
	log       logger.Logger
}

// New инициализирует приложение со всеми зависимостями
func New(cfg *config.Config) (*App, error) {
	// Инициализируем логгер
	log := logger.New("info")
	log.Info("Инициализация приложения")

	// Настройка Gin
	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	// Health check - независимо от других сервисов
	healthHandler := http.NewHealthHandler()
	router.GET("/api/v1/health", healthHandler.HealthCheck)

	// Подключение к PostgreSQL
	db, err := pgClient.NewConnection(cfg.Database)
	if err != nil {
		log.Warn("Не удалось подключиться к PostgreSQL, сервер запустится в ограниченном режиме", "error", err)
		return &App{router: router, log: log}, nil
	}
	log.Info("PostgreSQL подключен успешно")

	// Подключение к Redis (опционально)
	var redisClientInstance *redisLib.Client
	var cacheAdapter *redisAdapter.CacheAdapter
	redisClientInstance = redisClient.NewClient(cfg.Redis)
	if redisClientInstance != nil {
		cacheAdapter = redisAdapter.NewCacheAdapter(redisClientInstance, log)
		log.Info("Redis подключен успешно")
	} else {
		log.Warn("Redis недоступен, кэширование отключено")
	}

	// Подключение к RabbitMQ (опционально)
	var messagingAdapter *rabbitmq.MessagingAdapter
	if cfg.RabbitMQ.URL != "" {
		messagingAdapter, err = rabbitmq.NewMessagingAdapter(cfg.RabbitMQ.URL, log)
		if err != nil {
			log.Warn("Не удалось подключиться к RabbitMQ, разбивка задач может не работать", "error", err)
		} else {
			log.Info("RabbitMQ подключен успешно")
		}
	}

	// Инициализация адаптеров
	userAdapter := postgresql.NewUserAdapter(db, log)
	taskAdapter := postgresql.NewTaskAdapter(db, log)
	splitRequestAdapter := postgresql.NewSplitRequestAdapter(db, log)

	// Инициализация Use Cases для пользователей
	getUserUC := usecases.NewGetUserUseCase(userAdapter, log)
	updateUserUC := usecases.NewUpdateUserUseCase(userAdapter, log)
	loginUC := usecases.NewLoginUseCase(userAdapter, log)
	registerUC := usecases.NewRegisterUseCase(userAdapter, log)

	// Инициализация Use Cases для задач
	createTaskUC := usecases.NewCreateTaskUseCase(taskAdapter, log)
	getTaskUC := usecases.NewGetTaskUseCase(taskAdapter, log)
	getTasksUC := usecases.NewGetTasksUseCase(taskAdapter, log)
	updateTaskUC := usecases.NewUpdateTaskUseCase(taskAdapter, log)
	deleteTaskUC := usecases.NewDeleteTaskUseCase(taskAdapter, log)

	// Use Cases для разбивки задач (требуют messaging)
	var splitTaskUC *usecases.SplitTaskUseCase
	if messagingAdapter != nil {
		splitTaskUC = usecases.NewSplitTaskUseCase(taskAdapter, splitRequestAdapter, messagingAdapter, log)
	}
	getSplitStatusUC := usecases.NewGetSplitStatusUseCase(splitRequestAdapter, cacheAdapter, log)

	// Инициализация контроллеров
	userHandler := http.NewUserHandler(getUserUC, updateUserUC, log)
	taskHandler := http.NewTaskHandler(createTaskUC, getTaskUC, getTasksUC, updateTaskUC, deleteTaskUC, splitTaskUC, getSplitStatusUC, log)
	authHandler := http.NewAuthHandler(loginUC, registerUC, log)

	// Swagger документация
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Публичные маршруты
	public := router.Group("/api/v1")
	{
		public.POST("/auth/login", authHandler.Login)
		public.POST("/auth/register", authHandler.Register)
	}

	// Защищенные маршруты
	protected := router.Group("/api/v1")
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
		if splitTaskUC != nil {
			protected.POST("/split", taskHandler.SplitTask)
		}
		protected.GET("/split/:id/status", taskHandler.GetSplitStatus)
	}

	log.Info("Приложение успешно инициализировано")

	return &App{
		router:    router,
		db:        db,
		redis:     redisClientInstance,
		messaging: messagingAdapter,
		log:       log,
	}, nil
}

// Start запускает приложение
func (a *App) Start(host, port string) error {
	serverAddr := host + ":" + port
	a.log.Info("Запуск сервера", "address", serverAddr)
	a.log.Info("Health check доступен", "url", "http://"+serverAddr+"/api/v1/health")
	a.log.Info("Swagger документация", "url", "http://"+serverAddr+"/swagger/index.html")

	return a.router.Run(serverAddr)
}

// Shutdown корректно завершает работу приложения
func (a *App) Shutdown() error {
	a.log.Info("Завершение работы приложения")

	if a.messaging != nil {
		if err := a.messaging.Close(); err != nil {
			a.log.Error("Ошибка закрытия RabbitMQ", "error", err)
		}
	}

	if a.redis != nil {
		if err := a.redis.Close(); err != nil {
			a.log.Error("Ошибка закрытия Redis", "error", err)
		}
	}

	a.log.Info("Приложение завершено")
	return nil
}

