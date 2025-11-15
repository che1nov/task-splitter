# TaskSplitter

Приложение для разбивки задач на подзадачи с помощью ИИ. Когда у тебя большая задача, а ты не знаешь с чего начать — просто закинь её сюда, и получишь готовый план действий.

Стек: Go + PostgreSQL + React + RabbitMQ

## Что умеет

- Разбивает задачи на подзадачи через OpenAI/GigaChat
- Обычный CRUD для задач (создать, посмотреть, обновить, удалить)
- Авторизация через Keycloak (можно и без неё в demo-режиме)
- Роли пользователей (free/premium)
- Асинхронная обработка через RabbitMQ — задачи разбиваются в фоне
- Redis для кеша результатов
- React фронт на TypeScript
- Все в Docker

## Как это работает

```
Фронт (React) → API (Go) → Worker (фоновая обработка)
                    ↓
                PostgreSQL + Redis + RabbitMQ
```

По сути API принимает запросы, кидает их в очередь RabbitMQ, а Worker подхватывает и обрабатывает. Результат сохраняется в БД и кеше.

## Что нужно
- Docker и docker-compose (для локального запуска)
- Node.js 18+ (для фронта)
- OpenAI API ключ (если хочешь реальную разбивку задач)

## Быстрый запуск

```bash
# Клонируем
git clone https://github.com/che1nov/task-splitter.git
cd task-splitter

# Копируем .env
cp env.example .env

# Если есть OpenAI ключ - добавь его в .env
# OPENAI_API_KEY=sk-...

# Запускаем всё
make docker-up
```

Готово! Приложение доступно на:
- Фронт: http://localhost:3000
- API: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html
- RabbitMQ UI: http://localhost:15672 (guest/guest)
- Keycloak: http://localhost:8081 (admin/admin)

### Полезные команды

```bash
make help           # Показать все доступные команды
make docker-up      # Запустить всё в Docker
make docker-down    # Остановить все контейнеры
make docker-logs    # Посмотреть логи
make run-api        # Запустить только API
make run-worker     # Запустить только Worker
make test           # Запустить тесты
make fmt            # Отформатировать код
```

## Локальная разработка

Если хочешь разрабатывать без докера:

**Backend:**
```bash
# Поднимаем только инфру
docker-compose up -d postgres redis rabbitmq

# Запускаем API
go run cmd/server/main.go

# В другом терминале запускаем Worker
go run cmd/worker/main.go
```

**Frontend:**
```bash
cd web
npm install
npm start
```

## API

Swagger документация: http://localhost:8080/swagger/index.html

Основные endpoint'ы:

**Авторизация:**
- `POST /api/v1/auth/login` - логин
- `POST /api/v1/auth/register` - регистрация

**Задачи:**
- `GET /api/v1/tasks` - список задач
- `POST /api/v1/tasks` - создать задачу
- `GET /api/v1/tasks/:id` - получить задачу
- `PUT /api/v1/tasks/:id` - обновить
- `DELETE /api/v1/tasks/:id` - удалить

**Разбивка:**
- `POST /api/v1/split` - отправить задачу на разбивку
- `GET /api/v1/split/:id/status` - проверить статус

Все защищенные endpoint'ы требуют JWT токен в заголовке `Authorization: Bearer <token>`

## Архитектура проекта

Проект построен по принципам Clean Architecture. Вот структура:

```
internal/
├── domain/          # Бизнес-логика и сущности
│   ├── user.go
│   ├── task.go
│   ├── errors.go
│   └── ...
├── dto/             # Данные для передачи между слоями
├── usecases/        # Use Cases (каждый endpoint = отдельный use case)
│   ├── create_task_usecase.go
│   ├── split_task_usecase.go
│   └── interfaces.go
├── adapters/        # Адаптеры к внешним системам
│   ├── postgresql/
│   ├── redis/
│   └── rabbitmq/
├── controllers/     # HTTP контроллеры
│   └── http/
└── app/             # Dependency Injection
```

**Правила:**
- Domain не зависит ни от чего
- Use Cases содержат всю бизнес-логику
- Adapters работают с БД/Redis/RabbitMQ
- Controllers только принимают HTTP и отдают ответы

## Конфигурация

Всё настраивается через переменные окружения. См. `env.example`

Основное:
- `DATABASE_URL` - PostgreSQL (или отдельные DB_HOST, DB_PORT и т.д.)
- `REDIS_HOST`, `REDIS_PORT` - для кеша
- `RABBITMQ_URL` - для очередей
- `OPENAI_API_KEY` - для разбивки задач через GPT
- `GIGACHAT_*` - для разбивки через GigaChat (российская альтернатива)

## Тестирование

```bash
# Backend
make test                 # Все тесты
make test-coverage        # С отчетом о покрытии
make test-unit            # Только unit тесты

# Frontend
make web-test

# Линтеры
make lint                 # Запустить golangci-lint
make fmt                  # Отформатировать код
make vet                  # Проверить код на ошибки
```

## Деплой

**Docker Compose (проще всего):**
```bash
make prod               # Запустить в production режиме
```

**VPS (Ubuntu/Debian):**
```bash
# Устанавливаем Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Клонируем и запускаем
git clone https://github.com/che1nov/task-splitter.git
cd task-splitter
cp env.example .env
# Настраиваем .env под свои нужды
make docker-up
```

**Kubernetes:**
```bash
kubectl apply -f k8s/
```

## Безопасность

- JWT токены для авторизации
- Можно прикрутить Keycloak для управления пользователями
- CORS настроен
- Валидация всех входящих данных
- Rate limiting (TODO)

## Что дальше

Планы на будущее:
- [ ] Добавить поддержку разных языков
- [ ] Интеграцию с календарями (Google Calendar, Яндекс.Календарь)
- [ ] Мобилку
- [ ] Статистику и аналитику по задачам
- [ ] Интеграцию с Jira/Trello
- [ ] Улучшить разбивку задач с помощью fine-tuning моделей
