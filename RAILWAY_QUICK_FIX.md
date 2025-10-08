# 🚀 Railway Deployment - Quick Fix Guide

## ❌ Проблема: `failed to connect to host=localhost`

Эта ошибка возникает, когда приложение пытается подключиться к `localhost:5432`, но в Railway база данных находится на другом хосте.

## ✅ Решение:

### 1. Добавьте PostgreSQL в Railway
1. В Railway Dashboard → ваш проект
2. Нажмите **"+ New"** → **"Database"** → **"PostgreSQL"**
3. Railway автоматически создаст базу данных и установит переменную `DATABASE_URL`

### 2. Проверьте переменные окружения
В Settings → Variables должны быть:

```bash
# Обязательные
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# AI (выберите один)
GIGACHAT_CLIENT_ID=61e2f9ef-d364-4190-8d41-5458d59872bf
GIGACHAT_SCOPE=GIGACHAT_API_PERS
GIGACHAT_AUTH_KEY=NjFlMmY5ZWYtZDM2NC00MTkwLThkNDEtNTQ1OGQ1OTg3MmJmOjUxZjBiZmNiLTk4MjMtNDhiNy04ZWRmLWY5ZTY5MmJjYzAxYg==

# Frontend (замените на ваш домен)
REACT_APP_API_URL=https://your-app-name.railway.app/api/v1
```

### 3. НЕ устанавливайте эти переменные:
❌ `DB_HOST`  
❌ `DB_PORT`  
❌ `DB_USER`  
❌ `DB_PASSWORD`  
❌ `DB_NAME`  

Railway автоматически предоставляет `DATABASE_URL` с правильными значениями.

### 4. Проверьте логи
После деплоя проверьте логи Railway. Вы должны увидеть:

```
🔗 Using DATABASE_URL directly: postgres://postgres:***@containers-us-west-123.railway.app:6543/railway?sslmode=require
🔌 Connecting to database with DSN: postgres://postgres:***@containers-us-west-123.railway.app:6543/railway?sslmode=require
```

Если видите `⚠️ DATABASE_URL not found`, значит PostgreSQL плагин не добавлен.

## 🔍 Отладка:

### Проверьте health check:
```bash
curl https://your-app-name.railway.app/api/v1/health
```

### Проверьте подключение к базе данных:
В логах Railway должно быть:
- ✅ `🔗 Using DATABASE_URL directly`
- ✅ `Database migrations completed successfully`

Если есть ошибки:
- ❌ `⚠️ DATABASE_URL not found` → Добавьте PostgreSQL плагин
- ❌ `failed to connect to host=localhost` → Railway не обновил код, перезапустите деплой

## 🚀 Готово!

После этих шагов приложение должно успешно подключиться к базе данных Railway и работать корректно.
