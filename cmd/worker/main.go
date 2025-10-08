package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"task-splitter/config"
	"task-splitter/internal/models"
	"task-splitter/internal/repository"
	"task-splitter/pkg/gigachat"
	"task-splitter/pkg/messaging"
	"task-splitter/pkg/postgres"
	"task-splitter/pkg/redis"
	"time"

	redisLib "github.com/redis/go-redis/v9"
)

func main() {
	log.Println("Starting NLP Worker...")

	// Загрузка конфигурации
	cfg := config.Load()

	// Инициализация базы данных
	db, err := postgres.NewConnection(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Инициализация Redis
	redisClient := redis.NewClient(cfg.Redis)

	// Инициализация GigaChat клиента
	gigachatClient := gigachat.NewClient(cfg.GigaChat)

	// Инициализация репозиториев
	taskRepo := repository.NewTaskRepository(db)
	templateRepo := repository.NewTemplateRepository(db)
	splitRequestRepo := repository.NewTaskSplitRequestRepository(db)

	// Инициализация сервиса сообщений
	messagingService := messaging.NewRabbitMQService()
	defer messagingService.Close()

	// Создаем worker
	worker := NewNLPWorker(
		gigachatClient,
		redisClient,
		taskRepo,
		templateRepo,
		splitRequestRepo,
		messagingService,
	)

	// Запускаем обработку сообщений
	if err := worker.Start(); err != nil {
		log.Fatal("Failed to start worker:", err)
	}

	log.Println("NLP Worker started successfully")

	// Бесконечный цикл для поддержания работы worker'а
	for {
		time.Sleep(time.Hour)
	}
}

// NLPWorker представляет worker для обработки NLP задач
type NLPWorker struct {
	gigachatClient   *gigachat.Client
	redisClient      *redisLib.Client
	taskRepo         repository.TaskRepository
	templateRepo     repository.TemplateRepository
	splitRequestRepo repository.TaskSplitRequestRepository
	messagingService messaging.MessagingService
}

// NewNLPWorker создает новый NLP worker
func NewNLPWorker(
	gigachatClient *gigachat.Client,
	redisClient *redisLib.Client,
	taskRepo repository.TaskRepository,
	templateRepo repository.TemplateRepository,
	splitRequestRepo repository.TaskSplitRequestRepository,
	messagingService messaging.MessagingService,
) *NLPWorker {
	return &NLPWorker{
		gigachatClient:   gigachatClient,
		redisClient:      redisClient,
		taskRepo:         taskRepo,
		templateRepo:     templateRepo,
		splitRequestRepo: splitRequestRepo,
		messagingService: messagingService,
	}
}

// Start запускает worker
func (w *NLPWorker) Start() error {
	log.Println("Starting message consumption...")

	// Подписываемся на сообщения запросов на разбивку
	return w.messagingService.ConsumeSplitTaskRequests(w.handleSplitTaskRequest)
}

// handleSplitTaskRequest обрабатывает запрос на разбивку задачи
func (w *NLPWorker) handleSplitTaskRequest(message messaging.SplitTaskMessage) error {
	log.Printf("Processing split task request: %s", message.RequestID)

	ctx := context.Background()

	// Для демо создаем mock объект без сохранения в базу данных
	splitRequest := &models.TaskSplitRequest{
		ID:         message.RequestID,
		TaskID:     message.TaskID,
		UserID:     message.UserID,
		Text:       message.Text,
		TemplateID: message.TemplateID,
		Status:     "processing",
	}
	// В реальном приложении здесь должно быть:
	// if err := w.splitRequestRepo.Create(splitRequest); err != nil {
	//     log.Printf("Failed to create split request: %v", err)
	//     return err
	// }

	// Получаем шаблон, если указан
	var template *models.Template
	if message.TemplateID != nil {
		t, err := w.templateRepo.GetByID(*message.TemplateID)
		if err != nil {
			log.Printf("Failed to get template: %v", err)
		} else {
			template = t
		}
	}

	// Формируем промпт для GigaChat
	prompt := w.buildPrompt(message.Text, template)

	// Отправляем запрос в GigaChat
	result, err := w.gigachatClient.SplitTask(ctx, prompt)
	if err != nil {
		log.Printf("GigaChat request failed: %v", err)

		// Обновляем статус на failed
		splitRequest.Status = "failed"
		splitRequest.Error = err.Error()
		// w.splitRequestRepo.Update(splitRequest) // Закомментировано для демо

		// Отправляем ответ об ошибке
		response := messaging.SplitTaskResponse{
			RequestID: message.RequestID,
			TaskID:    message.TaskID,
			UserID:    message.UserID,
			Status:    "failed",
			Error:     err.Error(),
		}

		return w.messagingService.PublishSplitTaskResponse(response)
	}

	// Парсим результат и создаем подзадачи
	subtasks, err := w.parseSubtasks(result)
	if err != nil {
		log.Printf("Failed to parse subtasks: %v", err)

		splitRequest.Status = "failed"
		splitRequest.Error = err.Error()
		// w.splitRequestRepo.Update(splitRequest) // Закомментировано для демо

		response := messaging.SplitTaskResponse{
			RequestID: message.RequestID,
			TaskID:    message.TaskID,
			UserID:    message.UserID,
			Status:    "failed",
			Error:     err.Error(),
		}

		return w.messagingService.PublishSplitTaskResponse(response)
	}

	// Сохраняем подзадачи в базе данных
	if err := w.saveSubtasks(message.TaskID, subtasks); err != nil {
		log.Printf("Failed to save subtasks: %v", err)

		splitRequest.Status = "failed"
		splitRequest.Error = err.Error()
		// w.splitRequestRepo.Update(splitRequest) // Закомментировано для демо

		response := messaging.SplitTaskResponse{
			RequestID: message.RequestID,
			TaskID:    message.TaskID,
			UserID:    message.UserID,
			Status:    "failed",
			Error:     err.Error(),
		}

		return w.messagingService.PublishSplitTaskResponse(response)
	}

	// Сохраняем результат в Redis для быстрого доступа
	if err := w.redisClient.Set(ctx, "split_result:"+message.RequestID, result, 24*time.Hour).Err(); err != nil {
		log.Printf("Failed to cache result in Redis: %v", err)
	}

	// Обновляем статус на completed
	splitRequest.Status = "completed"
	splitRequest.Result = result
	// w.splitRequestRepo.Update(splitRequest) // Закомментировано для демо

	// Отправляем успешный ответ
	response := messaging.SplitTaskResponse{
		RequestID: message.RequestID,
		TaskID:    message.TaskID,
		UserID:    message.UserID,
		Status:    "completed",
		Result:    result,
	}

	log.Printf("Successfully processed split task request: %s", message.RequestID)
	return w.messagingService.PublishSplitTaskResponse(response)
}

// buildPrompt формирует промпт для GigaChat
func (w *NLPWorker) buildPrompt(text string, template *models.Template) string {
	basePrompt := `Ты эксперт по планированию и разбивке задач. Твоя задача - разбить большую задачу на максимально мелкие и конкретные шаги.

ПРАВИЛА РАЗБИВКИ:
1. Каждый шаг должен занимать 3-15 минут (максимум!)
2. Шаги должны быть настолько простыми, что их можно выполнить за один подход
3. Избегай абстрактных понятий - используй конкретные действия
4. Включай подготовительные шаги (найти материалы, настроить среду)
5. Добавляй проверочные шаги (проверить результат, протестировать)
6. Разбивай даже простые действия на микро-шаги

ПРИМЕРЫ ХОРОШИХ МИКРО-ШАГОВ:
❌ Плохо: "Изучить React" (слишком общее)
❌ Плохо: "Прочитать документацию" (слишком общее)
✅ Хорошо: "Открыть браузер", "Перейти на reactjs.org", "Найти раздел 'Getting Started'", "Прочитать первый абзац"

❌ Плохо: "Настроить проект" (неконкретно)  
✅ Хорошо: "Создать папку проекта", "Открыть терминал", "Выполнить команду 'npm init'", "Ответить на вопросы npm"

❌ Плохо: "Изучить useState" (слишком общее)
✅ Хорошо: "Найти пример useState в документации", "Прочитать синтаксис useState", "Скопировать пример кода", "Вставить в редактор"

ЦЕЛЬ: Создать 8-15 микро-шагов по 3-15 минут каждый.

Верни результат в формате JSON массив объектов с полями: title, description, priority (low/medium/high), estimated_time (в минутах).`

	if template != nil && template.Prompt != "" {
		basePrompt = template.Prompt
	}

	return fmt.Sprintf("%s\n\nЗадача: %s\n\nОтвет:", basePrompt, text)
}

// parseSubtasks парсит результат GigaChat в структуры подзадач
func (w *NLPWorker) parseSubtasks(result string) ([]models.Subtask, error) {
	var subtasks []models.Subtask

	// Пытаемся распарсить JSON
	var jsonSubtasks []map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonSubtasks); err != nil {
		// Если не JSON, пытаемся извлечь подзадачи из текста
		return w.extractSubtasksFromText(result)
	}

	// Конвертируем JSON в структуры
	for i, item := range jsonSubtasks {
		subtask := models.Subtask{
			Title:       getString(item, "title"),
			Description: getString(item, "description"),
			Priority:    getString(item, "priority"),
			Order:       i,
		}

		// Парсим estimated_time если есть
		if estimatedTime, ok := item["estimated_time"].(float64); ok {
			time := int(estimatedTime)
			subtask.EstimatedTime = &time
		}

		// Устанавливаем значения по умолчанию
		if subtask.Priority == "" {
			subtask.Priority = "medium"
		}

		subtasks = append(subtasks, subtask)
	}

	return subtasks, nil
}

// extractSubtasksFromText извлекает подзадачи из текстового ответа
func (w *NLPWorker) extractSubtasksFromText(text string) ([]models.Subtask, error) {
	// Простая реализация извлечения подзадач из текста
	// В реальном приложении здесь должна быть более сложная логика парсинга

	lines := strings.Split(text, "\n")
	var subtasks []models.Subtask

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "*") {
			continue
		}

		// Убираем номера и маркеры
		line = regexp.MustCompile(`^\d+\.\s*`).ReplaceAllString(line, "")
		line = regexp.MustCompile(`^[-*]\s*`).ReplaceAllString(line, "")

		if len(line) > 10 { // Минимальная длина подзадачи
			subtask := models.Subtask{
				Title:       line,
				Description: "",
				Priority:    "medium",
				Order:       len(subtasks),
			}
			subtasks = append(subtasks, subtask)
		}
	}

	return subtasks, nil
}

// saveSubtasks сохраняет подзадачи в базе данных
func (w *NLPWorker) saveSubtasks(taskID uint, subtasks []models.Subtask) error {
	for _, subtask := range subtasks {
		subtask.TaskID = taskID
		subtask.Status = "pending"

		// Здесь должен быть репозиторий для подзадач
		// Пока просто логируем
		log.Printf("Would save subtask: %s", subtask.Title)
	}

	return nil
}

// getString безопасно извлекает строку из map
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}
