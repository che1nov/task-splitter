package main

import (
	"log"
	"task-splitter/config"
	"task-splitter/internal/app"
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
	// Загружаем конфигурацию
	cfg := config.Load()

	// Инициализируем приложение (dependency injection)
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Ошибка инициализации приложения: %v", err)
	}

	// Запускаем HTTP сервер
	if err := application.Start(cfg.Server.Host, cfg.Server.Port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
