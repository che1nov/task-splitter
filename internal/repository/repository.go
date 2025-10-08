package repository

import (
	"task-splitter/internal/models"

	"gorm.io/gorm"
)

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByKeycloakID(keycloakID string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	List(limit, offset int) ([]*models.User, error)
}

// TaskRepository интерфейс для работы с задачами
type TaskRepository interface {
	Create(task *models.Task) error
	GetByID(id uint) (*models.Task, error)
	GetByUserID(userID uint, limit, offset int) ([]*models.Task, error)
	Update(task *models.Task) error
	Delete(id uint) error
	GetWithSubtasks(id uint) (*models.Task, error)
}

// TemplateRepository интерфейс для работы с шаблонами
type TemplateRepository interface {
	Create(template *models.Template) error
	GetByID(id uint) (*models.Template, error)
	GetByCategory(category string) ([]*models.Template, error)
	GetActive() ([]*models.Template, error)
	Update(template *models.Template) error
	Delete(id uint) error
	List(limit, offset int) ([]*models.Template, error)
}

// TaskSplitRequestRepository интерфейс для работы с запросами разбивки
type TaskSplitRequestRepository interface {
	Create(request *models.TaskSplitRequest) error
	GetByID(id string) (*models.TaskSplitRequest, error)
	GetByUserID(userID uint, limit, offset int) ([]*models.TaskSplitRequest, error)
	GetByStatus(status string) ([]*models.TaskSplitRequest, error)
	Update(request *models.TaskSplitRequest) error
	Delete(id string) error
}

// userRepository реализация UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository создает новый репозиторий пользователей
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByKeycloakID(keycloakID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("keycloak_id = ?", keycloakID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *userRepository) List(limit, offset int) ([]*models.User, error) {
	var users []*models.User
	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// taskRepository реализация TaskRepository
type taskRepository struct {
	db *gorm.DB
}

// NewTaskRepository создает новый репозиторий задач
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) GetByID(id uint) (*models.Task, error) {
	var task models.Task
	err := r.db.Preload("User").First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) GetByUserID(userID uint, limit, offset int) ([]*models.Task, error) {
	var tasks []*models.Task
	err := r.db.Where("user_id = ?", userID).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) Update(task *models.Task) error {
	return r.db.Save(task).Error
}

func (r *taskRepository) Delete(id uint) error {
	return r.db.Delete(&models.Task{}, id).Error
}

func (r *taskRepository) GetWithSubtasks(id uint) (*models.Task, error) {
	var task models.Task
	err := r.db.Preload("User").
		Preload("Subtasks", func(db *gorm.DB) *gorm.DB {
			return db.Order("order ASC")
		}).
		First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// templateRepository реализация TemplateRepository
type templateRepository struct {
	db *gorm.DB
}

// NewTemplateRepository создает новый репозиторий шаблонов
func NewTemplateRepository(db *gorm.DB) TemplateRepository {
	return &templateRepository{db: db}
}

func (r *templateRepository) Create(template *models.Template) error {
	return r.db.Create(template).Error
}

func (r *templateRepository) GetByID(id uint) (*models.Template, error) {
	var template models.Template
	err := r.db.First(&template, id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *templateRepository) GetByCategory(category string) ([]*models.Template, error) {
	var templates []*models.Template
	err := r.db.Where("category = ? AND is_active = ?", category, true).Find(&templates).Error
	return templates, err
}

func (r *templateRepository) GetActive() ([]*models.Template, error) {
	var templates []*models.Template
	err := r.db.Where("is_active = ?", true).Find(&templates).Error
	return templates, err
}

func (r *templateRepository) Update(template *models.Template) error {
	return r.db.Save(template).Error
}

func (r *templateRepository) Delete(id uint) error {
	return r.db.Delete(&models.Template{}, id).Error
}

func (r *templateRepository) List(limit, offset int) ([]*models.Template, error) {
	var templates []*models.Template
	err := r.db.Limit(limit).Offset(offset).Find(&templates).Error
	return templates, err
}

// taskSplitRequestRepository реализация TaskSplitRequestRepository
type taskSplitRequestRepository struct {
	db *gorm.DB
}

// NewTaskSplitRequestRepository создает новый репозиторий запросов разбивки
func NewTaskSplitRequestRepository(db *gorm.DB) TaskSplitRequestRepository {
	return &taskSplitRequestRepository{db: db}
}

func (r *taskSplitRequestRepository) Create(request *models.TaskSplitRequest) error {
	return r.db.Create(request).Error
}

func (r *taskSplitRequestRepository) GetByID(id string) (*models.TaskSplitRequest, error) {
	var request models.TaskSplitRequest
	err := r.db.Preload("Task").Preload("User").Preload("Template").
		First(&request, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *taskSplitRequestRepository) GetByUserID(userID uint, limit, offset int) ([]*models.TaskSplitRequest, error) {
	var requests []*models.TaskSplitRequest
	err := r.db.Where("user_id = ?", userID).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&requests).Error
	return requests, err
}

func (r *taskSplitRequestRepository) GetByStatus(status string) ([]*models.TaskSplitRequest, error) {
	var requests []*models.TaskSplitRequest
	err := r.db.Where("status = ?", status).Find(&requests).Error
	return requests, err
}

func (r *taskSplitRequestRepository) Update(request *models.TaskSplitRequest) error {
	return r.db.Save(request).Error
}

func (r *taskSplitRequestRepository) Delete(id string) error {
	return r.db.Delete(&models.TaskSplitRequest{}, "id = ?", id).Error
}
