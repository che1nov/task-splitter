# 🚀 TaskSplitter - Railway Deployment Guide

## 📋 Переменные окружения для Railway

### 🔧 Обязательные переменные:

```bash
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# Database (Railway автоматически предоставляет)
# DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME будут установлены автоматически
DB_SSLMODE=require

# AI Configuration (выберите один)
# GigaChat (рекомендуется)
GIGACHAT_CLIENT_ID=61e2f9ef-d364-4190-8d41-5458d59872bf
GIGACHAT_SCOPE=GIGACHAT_API_PERS
GIGACHAT_AUTH_KEY=NjFlMmY5ZWYtZDM2NC00MTkwLThkNDEtNTQ1OGQ1OTg3MmJmOjUxZjBiZmNiLTk4MjMtNDhiNy04ZWRmLWY5ZTY5MmJjYzAxYg==

# ИЛИ OpenAI
# OPENAI_API_KEY=your_openai_api_key_here
# OPENAI_MODEL=gpt-3.5-turbo

# Frontend URL (замените на ваш Railway домен)
REACT_APP_API_URL=https://your-app-name.railway.app/api/v1
```

### 🔗 Дополнительные сервисы (опционально):

```bash
# Redis (если подключаете)
REDIS_HOST=your-redis-host
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password
REDIS_DB=0

# RabbitMQ (если подключаете)
RABBITMQ_HOST=your-rabbitmq-host
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_VHOST=/
```

## 🚀 Пошаговая инструкция деплоя:

### 1. Подготовка репозитория
```bash
git add .
git commit -m "Add Railway deployment configuration"
git push origin main
```

### 2. Создание проекта в Railway
1. Зайдите на [railway.app](https://railway.app)
2. Нажмите "New Project"
3. Выберите "Deploy from GitHub repo"
4. Выберите ваш репозиторий `che1nov/task-splitter`

### 3. Настройка базы данных
1. В Railway Dashboard добавьте PostgreSQL плагин
2. Railway автоматически установит переменные:
   - `DATABASE_URL`
   - `DB_HOST`
   - `DB_PORT`
   - `DB_USER`
   - `DB_PASSWORD`
   - `DB_NAME`

### 4. Настройка переменных окружения
1. Перейдите в Settings → Variables
2. Добавьте все переменные из списка выше
3. **Важно:** Замените `REACT_APP_API_URL` на ваш реальный Railway домен

### 5. Деплой
1. Railway автоматически начнет сборку
2. Проверьте логи в разделе Deployments
3. После успешного деплоя получите URL приложения

## 🔍 Проверка работы:

### Health Check
```bash
curl https://your-app-name.railway.app/api/v1/health
```

### Регистрация пользователя
```bash
curl -X POST https://your-app-name.railway.app/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "firstName": "Test",
    "lastName": "User"
  }'
```

### Вход пользователя
```bash
curl -X POST https://your-app-name.railway.app/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

## 🛠️ Troubleshooting:

### Проблема: База данных не подключается
- Проверьте, что PostgreSQL плагин добавлен
- Убедитесь, что `DB_SSLMODE=require`

### Проблема: AI не работает
- Проверьте переменные `GIGACHAT_*` или `OPENAI_*`
- Убедитесь, что ключи действительны

### Проблема: Frontend не подключается к API
- Обновите `REACT_APP_API_URL` на правильный Railway домен
- Проверьте CORS настройки

## 📊 Мониторинг:

Railway предоставляет:
- Логи приложения в реальном времени
- Метрики использования
- Автоматические перезапуски при сбоях
- Health check на `/api/v1/health`

## 🔒 Безопасность:

- Все секретные ключи хранятся в Railway Variables
- База данных изолирована
- HTTPS включен по умолчанию
- Автоматические обновления безопасности
