package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"task-splitter/internal/models"
	"task-splitter/internal/repository"
	"task-splitter/pkg/keycloak"
	"task-splitter/pkg/messaging"
	"time"

	redisLib "github.com/redis/go-redis/v9"
)

// UserService интерфейс для работы с пользователями
type UserService interface {
	CreateUser(req *models.CreateUserRequest) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetUserByKeycloakID(keycloakID string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUser(id uint, req *models.UpdateUserRequest) (*models.User, error)
	DeleteUser(id uint) error
	ListUsers(limit, offset int) ([]*models.User, error)
}

// TaskService интерфейс для работы с задачами
type TaskService interface {
	CreateTask(userID uint, req *models.CreateTaskRequest) (*models.Task, error)
	GetTaskByID(id uint) (*models.Task, error)
	GetTasksByUserID(userID uint, limit, offset int) ([]*models.Task, error)
	UpdateTask(id uint, req *models.UpdateTaskRequest) (*models.Task, error)
	DeleteTask(id uint) error
	SplitTask(userID uint, req *models.SplitTaskRequest) (*models.SplitTaskResponse, error)
	GetSplitStatus(requestID string) (*models.SplitStatusResponse, error)
}

// AuthService интерфейс для аутентификации
type AuthService interface {
	Login(ctx context.Context, username, password string) (*keycloak.UserInfo, error)
	Register(ctx context.Context, username, email, password, firstName, lastName string) (*models.User, error)
	ValidateToken(ctx context.Context, token string) (*keycloak.TokenInfo, error)
}

// userService реализация UserService
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService создает новый сервис пользователей
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	user := &models.User{
		KeycloakID: req.KeycloakID,
		Email:      req.Email,
		Username:   req.Username,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Role:       req.Role,
	}

	if user.Role == "" {
		user.Role = "free"
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) GetUserByKeycloakID(keycloakID string) (*models.User, error) {
	return s.userRepo.GetByKeycloakID(keycloakID)
}

func (s *userService) GetUserByUsername(username string) (*models.User, error) {
	return s.userRepo.GetByUsername(username)
}

func (s *userService) UpdateUser(id uint, req *models.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Role != "" {
		user.Role = req.Role
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *userService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *userService) ListUsers(limit, offset int) ([]*models.User, error) {
	return s.userRepo.List(limit, offset)
}

// taskService реализация TaskService
type taskService struct {
	taskRepo         repository.TaskRepository
	templateRepo     repository.TemplateRepository
	splitRequestRepo repository.TaskSplitRequestRepository
	redisClient      *redisLib.Client
	messaging        messaging.MessagingService
}

// NewTaskService создает новый сервис задач
func NewTaskService(taskRepo repository.TaskRepository, templateRepo repository.TemplateRepository, splitRequestRepo repository.TaskSplitRequestRepository, redisClient *redisLib.Client) TaskService {
	return &taskService{
		taskRepo:         taskRepo,
		templateRepo:     templateRepo,
		splitRequestRepo: splitRequestRepo,
		redisClient:      redisClient,
		messaging:        messaging.NewRabbitMQService(),
	}
}

func (s *taskService) CreateTask(userID uint, req *models.CreateTaskRequest) (*models.Task, error) {
	task := &models.Task{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		Status:      "pending",
	}

	if task.Priority == "" {
		task.Priority = "medium"
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

func (s *taskService) GetTaskByID(id uint) (*models.Task, error) {
	return s.taskRepo.GetWithSubtasks(id)
}

func (s *taskService) GetTasksByUserID(userID uint, limit, offset int) ([]*models.Task, error) {
	return s.taskRepo.GetByUserID(userID, limit, offset)
}

func (s *taskService) UpdateTask(id uint, req *models.UpdateTaskRequest) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Status != "" {
		task.Status = req.Status
	}
	if req.Priority != "" {
		task.Priority = req.Priority
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	if err := s.taskRepo.Update(task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return task, nil
}

func (s *taskService) DeleteTask(id uint) error {
	return s.taskRepo.Delete(id)
}

func (s *taskService) SplitTask(userID uint, req *models.SplitTaskRequest) (*models.SplitTaskResponse, error) {
	// Проверяем существование задачи
	task, err := s.taskRepo.GetByID(req.TaskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}
	if task.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	// Генерируем уникальный ID запроса
	requestID, err := generateRequestID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate request ID: %w", err)
	}

	// Создаем запрос на разбивку
	splitRequest := &models.TaskSplitRequest{
		ID:         requestID,
		TaskID:     req.TaskID,
		UserID:     userID,
		Text:       req.Text,
		TemplateID: req.TemplateID,
		Status:     "pending",
	}

	// Отправляем запрос в очередь
	message := messaging.SplitTaskMessage{
		RequestID:  requestID,
		TaskID:     req.TaskID,
		UserID:     userID,
		Text:       req.Text,
		TemplateID: req.TemplateID,
	}

	if err := s.messaging.PublishSplitTaskRequest(message); err != nil {
		return nil, fmt.Errorf("failed to publish split task request: %w", err)
	}

	// Сохраняем запрос в базе данных
	// В реальном приложении это должно быть в транзакции
	if err := s.splitRequestRepo.Create(splitRequest); err != nil {
		log.Printf("Failed to save split request: %v", err)
		// Не возвращаем ошибку, так как сообщение уже отправлено
	}
	log.Printf("Split task request created: %s", requestID)

	return &models.SplitTaskResponse{
		RequestID: requestID,
		Status:    "pending",
		Message:   "Task split request submitted successfully",
	}, nil
}

func (s *taskService) GetSplitStatus(requestID string) (*models.SplitStatusResponse, error) {
	// Сначала проверяем кеш Redis
	ctx := context.Background()
	cachedResult, err := s.redisClient.Get(ctx, "split_result:"+requestID).Result()
	if err == nil {
		// Результат найден в кеше
		return &models.SplitStatusResponse{
			RequestID: requestID,
			Status:    "completed",
			Result:    cachedResult,
		}, nil
	}

	// Если не в кеше, проверяем базу данных
	// В реальном приложении здесь должен быть репозиторий для TaskSplitRequest
	return &models.SplitStatusResponse{
		RequestID: requestID,
		Status:    "processing",
	}, nil
}

// authService реализация AuthService
type authService struct {
	keycloakClient *keycloak.Client
	userService    UserService
}

// NewAuthService создает новый сервис аутентификации
func NewAuthService(keycloakClient *keycloak.Client, userService UserService) AuthService {
	return &authService{
		keycloakClient: keycloakClient,
		userService:    userService,
	}
}

func (s *authService) Login(ctx context.Context, username, password string) (*keycloak.UserInfo, error) {
	// Ищем пользователя по username в базе данных
	user, err := s.userService.GetUserByUsername(username)
	if err != nil {
		// Если пользователь не найден, возвращаем ошибку
		return nil, fmt.Errorf("пользователь не найден")
	}

	// В реальном приложении здесь должна быть проверка пароля
	// Пока что просто проверяем, что пользователь существует

	return &keycloak.UserInfo{
		Sub:               user.KeycloakID,
		Email:             user.Email,
		PreferredUsername: user.Username,
		GivenName:         user.FirstName,
		FamilyName:        user.LastName,
		RealmRoles:        []string{user.Role},
		Groups:            []string{},
	}, nil
}

func (s *authService) Register(ctx context.Context, username, email, password, firstName, lastName string) (*models.User, error) {
	// Генерируем уникальный KeycloakID для демо-режима
	keycloakID := "demo_" + username + "_" + fmt.Sprintf("%d", time.Now().Unix())

	// Создаем пользователя через UserService
	createReq := &models.CreateUserRequest{
		KeycloakID: keycloakID,
		Email:      email,
		Username:   username,
		FirstName:  firstName,
		LastName:   lastName,
		Role:       "free",
	}

	user, err := s.userService.CreateUser(createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*keycloak.TokenInfo, error) {
	return s.keycloakClient.ValidateToken(ctx, token)
}

// generateRequestID генерирует уникальный ID запроса
func generateRequestID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
