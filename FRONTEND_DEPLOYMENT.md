# Развертывание фронтенда на Railway

## Пошаговая инструкция

### 1. Создание нового сервиса в Railway

1. Откройте панель Railway (https://railway.app)
2. Выберите ваш проект `task-splitter`
3. Нажмите **"New Service"** → **"GitHub Repo"**
4. Выберите репозиторий `che1nov/task-splitter`

### 2. Настройка сервиса

**Вариант A: Использование Dockerfile (рекомендуется)**

1. В настройках сервиса укажите:
   - **Dockerfile Path**: `Dockerfile.frontend`
   - Railway автоматически найдет и использует этот файл
   - Никаких дополнительных команд не требуется (nginx запускается автоматически)

**Вариант B: Использование Node.js**

1. В настройках сервиса укажите:
   - **Root Directory**: `web`
   - **Build Command**: `npm run build`
   - **Start Command**: `npm start`

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
├── Dockerfile.frontend # Docker конфигурация для фронтенда (в корне)
├── web/
│   ├── Dockerfile      # Docker конфигурация для локальной разработки
│   ├── railway.json    # Конфигурация Railway
│   ├── railway.env      # Переменные окружения для Railway
│   ├── nginx.conf       # Конфигурация nginx
│   ├── package.json     # Зависимости и скрипты
│   └── src/
│       └── App.tsx      # Обновлен для работы с переменными окружения
```

## Примечания

- Фронтенд будет доступен по отдельному домену
- API URL настраивается через переменную `REACT_APP_API_URL`
- Используется nginx для обслуживания статических файлов
- Поддерживается сжатие и кеширование
