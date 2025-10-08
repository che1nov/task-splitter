# TaskSplitter

TaskSplitter — это современное веб-приложение для разбивки задач на подзадачи с использованием искусственного интеллекта. Приложение построено на стеке Go, PostgreSQL, Keycloak и React.

## 🚀 Возможности

- **ИИ-разбивка задач**: Автоматическое разбиение задач на логические подзадачи с помощью OpenAI GPT-4
- **Управление задачами**: Полный CRUD для задач и подзадач
- **Аутентификация**: Интеграция с Keycloak для безопасной аутентификации
- **Роли пользователей**: Поддержка ролей free и premium
- **Асинхронная обработка**: Использование RabbitMQ для обработки запросов на разбивку
- **Кеширование**: Redis для быстрого доступа к результатам
- **Современный UI**: React-интерфейс с TypeScript
- **Docker**: Полная контейнеризация всех сервисов

## 🏗️ Архитектура

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   React UI      │    │   Go API        │    │   NLP Worker    │
│   (Port 3000)   │◄──►│   (Port 8080)   │◄──►│   (Background)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Keycloak      │    │   PostgreSQL    │    │   Redis         │
│   (Port 8081)   │    │   (Port 5432)   │    │   (Port 6379)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   RabbitMQ      │
                       │   (Port 5672)   │
                       └─────────────────┘
```

## 📋 Требования

- Docker и Docker Compose
- Go 1.21+ (для локальной разработки)
- Node.js 18+ (для frontend разработки)
- OpenAI API ключ (опционально, для ИИ-функций)

## 🚀 Быстрый старт

### 1. Клонирование репозитория

```bash
git clone https://github.com/your-username/task-splitter.git
cd task-splitter
```

### 2. Настройка переменных окружения

```bash
cp env.example .env
```

Отредактируйте `.env` файл и добавьте ваш OpenAI API ключ:

```env
OPENAI_API_KEY=your_openai_api_key_here
```

### 3. Запуск с Docker Compose

```bash
docker-compose up -d
```

Это запустит все сервисы:
- **API**: http://localhost:8080
- **Frontend**: http://localhost:3000
- **Keycloak**: http://localhost:8081
- **RabbitMQ Management**: http://localhost:15672

### 4. Инициализация Keycloak

1. Откройте http://localhost:8081
2. Войдите как admin/admin
3. Создайте realm "tasksplitter"
4. Создайте клиента "tasksplitter-api"
5. Создайте пользователей и назначьте роли

### 5. Проверка работы

Откройте http://localhost:3000 и протестируйте приложение.

## 🔧 Локальная разработка

### Backend (Go)

```bash
# Установка зависимостей
go mod download

# Запуск базы данных
docker-compose up -d postgres redis rabbitmq keycloak

# Запуск API сервера
go run cmd/server/main.go

# Запуск Worker
go run cmd/worker/main.go
```

### Frontend (React)

```bash
cd web

# Установка зависимостей
npm install

# Запуск в режиме разработки
npm start
```

## 📚 API Документация

После запуска API сервера, Swagger документация доступна по адресу:
http://localhost:8080/swagger/index.html

### Основные эндпоинты

#### Аутентификация
- `POST /api/v1/auth/login` - Вход пользователя
- `POST /api/v1/auth/register` - Регистрация пользователя

#### Задачи
- `GET /api/v1/tasks` - Получить список задач
- `POST /api/v1/tasks` - Создать задачу
- `GET /api/v1/tasks/{id}` - Получить задачу по ID
- `PUT /api/v1/tasks/{id}` - Обновить задачу
- `DELETE /api/v1/tasks/{id}` - Удалить задачу

#### Разбивка задач
- `POST /api/v1/split` - Отправить задачу на разбивку
- `GET /api/v1/split/{id}/status` - Получить статус разбивки

## 🧪 Тестирование

### Backend тесты

```bash
go test -v ./...
```

### Frontend тесты

```bash
cd web
npm test
```

### Интеграционные тесты

```bash
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

## 📦 Развертывание

### Production развертывание

1. Настройте переменные окружения для production
2. Используйте production Docker Compose файл:

```bash
docker-compose -f docker-compose.prod.yml up -d
```

### Kubernetes

```bash
kubectl apply -f k8s/
```

## 🔒 Безопасность

- Все API эндпоинты защищены JWT токенами
- Интеграция с Keycloak для управления пользователями
- CORS настроен для безопасности
- Валидация входных данных
- Rate limiting для предотвращения злоупотреблений

## 📊 Мониторинг

- Логирование всех операций
- Метрики производительности
- Health checks для всех сервисов
- Prometheus метрики (планируется)

## 🤝 Вклад в проект

1. Fork репозитория
2. Создайте feature branch (`git checkout -b feature/amazing-feature`)
3. Commit изменения (`git commit -m 'Add amazing feature'`)
4. Push в branch (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

## 📄 Лицензия

Этот проект лицензирован под MIT License - см. файл [LICENSE](LICENSE) для деталей.

## 🆘 Поддержка

Если у вас есть вопросы или проблемы:

1. Проверьте [Issues](https://github.com/your-username/task-splitter/issues)
2. Создайте новый Issue с подробным описанием
3. Свяжитесь с командой разработки

## 🗺️ Roadmap

- [ ] Поддержка дополнительных языков
- [ ] Интеграция с календарями
- [ ] Мобильное приложение
- [ ] Расширенная аналитика
- [ ] Интеграция с внешними сервисами (Jira, Trello)
- [ ] Машинное обучение для улучшения разбивки задач
