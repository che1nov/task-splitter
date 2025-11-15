.PHONY: help run-api run-worker run dev build test clean docker-up docker-down docker-logs migrate lint fmt deps

# Цвета для вывода
GREEN  := \033[0;32m
YELLOW := \033[0;33m
NC     := \033[0m

help: ## Показать это сообщение
	@echo "$(GREEN)TaskSplitter - Доступные команды:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

# Запуск приложения
run-api: ## Запустить API сервер
	@echo "$(GREEN)Запускаем API сервер...$(NC)"
	go run cmd/server/main.go

run-worker: ## Запустить Worker
	@echo "$(GREEN)Запускаем Worker...$(NC)"
	go run cmd/worker/main.go

run: ## Запустить API и Worker одновременно
	@echo "$(GREEN)Запускаем API и Worker...$(NC)"
	@make -j2 run-api run-worker

dev: ## Запустить в режиме разработки (с hot reload)
	@echo "$(GREEN)Запускаем в режиме разработки...$(NC)"
	air

# Сборка
build: ## Собрать все бинарники
	@echo "$(GREEN)Собираем приложение...$(NC)"
	@mkdir -p bin
	go build -o bin/server cmd/server/main.go
	go build -o bin/worker cmd/worker/main.go
	@echo "$(GREEN)✓ Бинарники собраны в ./bin/$(NC)"

build-docker: ## Собрать Docker образы
	@echo "$(GREEN)Собираем Docker образы...$(NC)"
	docker-compose build

# Тестирование
test: ## Запустить все тесты
	@echo "$(GREEN)Запускаем тесты...$(NC)"
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Запустить тесты с отчетом о покрытии
	@echo "$(GREEN)Генерируем отчет о покрытии...$(NC)"
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Отчет сохранен в coverage.html$(NC)"

test-unit: ## Запустить только unit тесты
	@echo "$(GREEN)Запускаем unit тесты...$(NC)"
	go test -v -short ./...

# Docker команды
docker-up: ## Запустить все сервисы в Docker
	@echo "$(GREEN)Запускаем Docker контейнеры...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)✓ Все сервисы запущены$(NC)"
	@echo "API: http://localhost:8080"
	@echo "Frontend: http://localhost:3000"
	@echo "Swagger: http://localhost:8080/swagger/index.html"

docker-down: ## Остановить все сервисы
	@echo "$(YELLOW)Останавливаем Docker контейнеры...$(NC)"
	docker-compose down

docker-restart: ## Перезапустить все сервисы
	@make docker-down
	@make docker-up

docker-logs: ## Показать логи всех сервисов
	docker-compose logs -f

docker-logs-api: ## Показать логи API
	docker-compose logs -f api

docker-logs-worker: ## Показать логи Worker
	docker-compose logs -f worker

docker-clean: ## Удалить все контейнеры и volumes
	@echo "$(YELLOW)Удаляем контейнеры и volumes...$(NC)"
	docker-compose down -v
	docker system prune -f

# База данных
db-up: ## Запустить только PostgreSQL
	@echo "$(GREEN)Запускаем PostgreSQL...$(NC)"
	docker-compose up -d postgres

migrate-up: ## Применить миграции
	@echo "$(GREEN)Применяем миграции...$(NC)"
	go run cmd/migrate/main.go up

migrate-down: ## Откатить миграции
	@echo "$(YELLOW)Откатываем миграции...$(NC)"
	go run cmd/migrate/main.go down

migrate-create: ## Создать новую миграцию (использование: make migrate-create name=имя_миграции)
	@echo "$(GREEN)Создаем миграцию...$(NC)"
	migrate create -ext sql -dir migrations -seq $(name)

# Инфраструктура
infra-up: ## Запустить только инфраструктуру (БД, Redis, RabbitMQ)
	@echo "$(GREEN)Запускаем инфраструктуру...$(NC)"
	docker-compose up -d postgres redis rabbitmq keycloak

infra-down: ## Остановить инфраструктуру
	docker-compose stop postgres redis rabbitmq keycloak

# Код
fmt: ## Отформатировать код
	@echo "$(GREEN)Форматируем код...$(NC)"
	go fmt ./...
	goimports -w .

lint: ## Запустить линтеры
	@echo "$(GREEN)Запускаем линтеры...$(NC)"
	golangci-lint run ./...

vet: ## Проверить код на ошибки
	@echo "$(GREEN)Проверяем код...$(NC)"
	go vet ./...

# Зависимости
deps: ## Установить зависимости
	@echo "$(GREEN)Устанавливаем зависимости...$(NC)"
	go mod download
	go mod tidy

deps-update: ## Обновить зависимости
	@echo "$(GREEN)Обновляем зависимости...$(NC)"
	go get -u ./...
	go mod tidy

# Frontend
web-install: ## Установить зависимости frontend
	@echo "$(GREEN)Устанавливаем npm пакеты...$(NC)"
	cd web && npm install

web-dev: ## Запустить frontend в режиме разработки
	@echo "$(GREEN)Запускаем frontend...$(NC)"
	cd web && npm start

web-build: ## Собрать frontend для продакшена
	@echo "$(GREEN)Собираем frontend...$(NC)"
	cd web && npm run build

web-test: ## Запустить тесты frontend
	cd web && npm test

# Swagger
swagger: ## Сгенерировать Swagger документацию
	@echo "$(GREEN)Генерируем Swagger документацию...$(NC)"
	swag init -g cmd/server/main.go -o ./docs

# Утилиты
clean: ## Очистить временные файлы
	@echo "$(YELLOW)Очищаем временные файлы...$(NC)"
	rm -rf bin/
	rm -rf coverage.out coverage.html
	rm -rf tmp/
	go clean

install-tools: ## Установить необходимые инструменты
	@echo "$(GREEN)Устанавливаем инструменты для разработки...$(NC)"
	go install github.com/cosmtrek/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "$(GREEN)✓ Все инструменты установлены$(NC)"

health: ## Проверить здоровье приложения
	@echo "$(GREEN)Проверяем health check...$(NC)"
	@curl -s http://localhost:8080/api/v1/health | jq . || echo "API не запущен"

# Полезные команды
all: clean deps build test ## Сборка с нуля (clean + deps + build + test)

prod: ## Запустить в production режиме
	@echo "$(GREEN)Запускаем в production режиме...$(NC)"
	docker-compose -f docker-compose.prod.yml up -d

setup: install-tools deps ## Первоначальная настройка проекта
	@echo "$(GREEN)Проект настроен! Запусти 'make docker-up' чтобы начать$(NC)"

