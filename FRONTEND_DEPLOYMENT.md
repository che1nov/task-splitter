# Развертывание фронтенда на Railway

## Пошаговая инструкция

### 1. Создание нового сервиса в Railway

1. Откройте панель Railway (https://railway.app)
2. Выберите ваш проект `task-splitter`
3. Нажмите **"New Service"** → **"GitHub Repo"**
4. Выберите репозиторий `che1nov/task-splitter`

### 2. Настройка сервиса

1. В настройках сервиса укажите:
   - **Root Directory**: `web`
   - **Build Command**: `npm run build`
   - **Start Command**: `npm start`

2. Или используйте Dockerfile (рекомендуется):
   - Railway автоматически найдет `Dockerfile` в папке `web/`
   - Никаких дополнительных настроек не требуется

### 3. Переменные окружения

Добавьте следующие переменные в настройках сервиса:

```
REACT_APP_API_URL=https://task-splitter-production.up.railway.app
NODE_ENV=production
GENERATE_SOURCEMAP=false
PORT=3000
```

### 4. Домен

После развертывания Railway автоматически назначит домен вида:
`https://your-frontend-service-name.up.railway.app`

### 5. Проверка

После развертывания проверьте:
- Фронтенд доступен по назначенному домену
- Фронтенд корректно подключается к API
- Все функции работают

## Альтернативный способ - через Railway CLI

```bash
# Установите Railway CLI
npm install -g @railway/cli

# Войдите в аккаунт
railway login

# Создайте новый сервис
railway service create frontend

# Установите переменные окружения
railway variables set REACT_APP_API_URL=https://task-splitter-production.up.railway.app
railway variables set NODE_ENV=production

# Разверните сервис
railway up --service frontend
```

## Структура файлов

```
web/
├── Dockerfile          # Docker конфигурация для Railway
├── railway.json        # Конфигурация Railway
├── railway.env         # Переменные окружения для Railway
├── nginx.conf          # Конфигурация nginx
├── package.json        # Зависимости и скрипты
└── src/
    └── App.tsx         # Обновлен для работы с переменными окружения
```

## Примечания

- Фронтенд будет доступен по отдельному домену
- API URL настраивается через переменную `REACT_APP_API_URL`
- Используется nginx для обслуживания статических файлов
- Поддерживается сжатие и кеширование
