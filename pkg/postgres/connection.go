package postgres

import (
	"fmt"
	"log"
	"task-splitter/config"
	"task-splitter/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewConnection создает новое подключение к PostgreSQL
func NewConnection(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Автоматическая миграция
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

// migrate выполняет миграции базы данных
func migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Миграция всех моделей
	err := db.AutoMigrate(
		&models.User{},
		&models.Task{},
		&models.Subtask{},
		&models.Template{},
		&models.TaskSplitRequest{},
	)
	if err != nil {
		return err
	}

	// Создание индексов
	if err := createIndexes(db); err != nil {
		return err
	}

	// Создание начальных данных
	if err := seedData(db); err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// createIndexes создает дополнительные индексы
func createIndexes(db *gorm.DB) error {
	// Индекс для поиска задач по пользователю
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id)").Error; err != nil {
		return err
	}

	// Индекс для поиска подзадач по задаче
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_subtasks_task_id ON subtasks(task_id)").Error; err != nil {
		return err
	}

	// Индекс для поиска запросов разбивки по статусу
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_task_split_requests_status ON task_split_requests(status)").Error; err != nil {
		return err
	}

	return nil
}

// seedData создает начальные данные
func seedData(db *gorm.DB) error {
	// Проверяем, есть ли уже шаблоны
	var count int64
	db.Model(&models.Template{}).Count(&count)
	if count > 0 {
		return nil // Данные уже есть
	}

	// Создаем базовые шаблоны
	templates := []models.Template{
		{
			Name:        "Общая разбивка задач",
			Description: "Универсальный шаблон для разбивки любых задач",
			Category:    "general",
			Prompt:      "Разбей следующую задачу на логические подзадачи с учетом приоритетов и зависимостей. Каждая подзадача должна быть конкретной и выполнимой.",
			IsActive:    true,
		},
		{
			Name:        "Разработка ПО",
			Description: "Шаблон для разбивки задач разработки программного обеспечения",
			Category:    "development",
			Prompt:      "Разбей задачу разработки на этапы: планирование, проектирование, реализация, тестирование, документация. Учти технические требования и зависимости.",
			IsActive:    true,
		},
		{
			Name:        "Маркетинговая кампания",
			Description: "Шаблон для разбивки маркетинговых задач",
			Category:    "marketing",
			Prompt:      "Разбей маркетинговую задачу на этапы: исследование, стратегия, контент, продвижение, анализ результатов. Учти целевую аудиторию и каналы продвижения.",
			IsActive:    true,
		},
		{
			Name:        "Учебный проект",
			Description: "Шаблон для разбивки учебных задач и проектов",
			Category:    "education",
			Prompt:      "Разбей учебную задачу на этапы: изучение теории, практические упражнения, проект, тестирование знаний. Учти уровень сложности и временные рамки.",
			IsActive:    true,
		},
	}

	for _, template := range templates {
		if err := db.Create(&template).Error; err != nil {
			return fmt.Errorf("failed to create template %s: %w", template.Name, err)
		}
	}

	log.Println("Initial data seeded successfully")
	return nil
}
