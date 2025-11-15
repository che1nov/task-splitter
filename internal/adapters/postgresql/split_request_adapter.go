package postgresql

import (
	"context"
	"task-splitter/internal/domain"
	"task-splitter/pkg/logger"

	"gorm.io/gorm"
)

// SplitRequestAdapter адаптер для работы с запросами разбивки в PostgreSQL
type SplitRequestAdapter struct {
	db  *gorm.DB
	log logger.Logger
}

// NewSplitRequestAdapter создает новый адаптер запросов разбивки
func NewSplitRequestAdapter(db *gorm.DB, log logger.Logger) *SplitRequestAdapter {
	return &SplitRequestAdapter{
		db:  db,
		log: log,
	}
}

// CreateSplitRequest создает запрос на разбивку в БД
func (a *SplitRequestAdapter) CreateSplitRequest(ctx context.Context, request domain.TaskSplitRequest) error {
	model := TaskSplitRequestModel{
		ID:         request.ID,
		TaskID:     request.TaskID,
		UserID:     request.UserID,
		Text:       request.Text,
		TemplateID: request.TemplateID,
		Status:     request.Status,
		Result:     request.Result,
		Error:      request.Error,
		CreatedAt:  request.CreatedAt,
		UpdatedAt:  request.UpdatedAt,
	}

	if err := a.db.WithContext(ctx).Create(&model).Error; err != nil {
		a.log.ErrorContext(ctx, "Ошибка создания запроса разбивки в БД", "error", err)
		return err
	}

	return nil
}

// GetSplitRequestByID получает запрос на разбивку по ID
func (a *SplitRequestAdapter) GetSplitRequestByID(ctx context.Context, id string) (domain.TaskSplitRequest, error) {
	var model TaskSplitRequestModel
	if err := a.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.TaskSplitRequest{}, domain.ErrSplitRequestNotFound
		}
		a.log.ErrorContext(ctx, "Ошибка получения запроса разбивки по ID", "id", id, "error", err)
		return domain.TaskSplitRequest{}, err
	}

	return toDomainSplitRequest(model), nil
}

// GetSplitRequestsByUserID получает запросы разбивки пользователя
func (a *SplitRequestAdapter) GetSplitRequestsByUserID(ctx context.Context, userID uint, limit, offset int) ([]domain.TaskSplitRequest, error) {
	var models []TaskSplitRequestModel
	if err := a.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		a.log.ErrorContext(ctx, "Ошибка получения запросов разбивки пользователя", "user_id", userID, "error", err)
		return nil, err
	}

	requests := make([]domain.TaskSplitRequest, len(models))
	for i, model := range models {
		requests[i] = toDomainSplitRequest(model)
	}

	return requests, nil
}

// GetSplitRequestsByStatus получает запросы разбивки по статусу
func (a *SplitRequestAdapter) GetSplitRequestsByStatus(ctx context.Context, status string) ([]domain.TaskSplitRequest, error) {
	var models []TaskSplitRequestModel
	if err := a.db.WithContext(ctx).Where("status = ?", status).Find(&models).Error; err != nil {
		a.log.ErrorContext(ctx, "Ошибка получения запросов разбивки по статусу", "status", status, "error", err)
		return nil, err
	}

	requests := make([]domain.TaskSplitRequest, len(models))
	for i, model := range models {
		requests[i] = toDomainSplitRequest(model)
	}

	return requests, nil
}

// UpdateSplitRequest обновляет запрос на разбивку
func (a *SplitRequestAdapter) UpdateSplitRequest(ctx context.Context, request domain.TaskSplitRequest) error {
	model := TaskSplitRequestModel{
		ID:         request.ID,
		TaskID:     request.TaskID,
		UserID:     request.UserID,
		Text:       request.Text,
		TemplateID: request.TemplateID,
		Status:     request.Status,
		Result:     request.Result,
		Error:      request.Error,
		CreatedAt:  request.CreatedAt,
		UpdatedAt:  request.UpdatedAt,
	}

	if err := a.db.WithContext(ctx).Save(&model).Error; err != nil {
		a.log.ErrorContext(ctx, "Ошибка обновления запроса разбивки", "id", request.ID, "error", err)
		return err
	}

	return nil
}

// DeleteSplitRequest удаляет запрос на разбивку
func (a *SplitRequestAdapter) DeleteSplitRequest(ctx context.Context, id string) error {
	if err := a.db.WithContext(ctx).Delete(&TaskSplitRequestModel{}, "id = ?", id).Error; err != nil {
		a.log.ErrorContext(ctx, "Ошибка удаления запроса разбивки", "id", id, "error", err)
		return err
	}

	return nil
}

// toDomainSplitRequest конвертирует модель БД в доменную модель
func toDomainSplitRequest(model TaskSplitRequestModel) domain.TaskSplitRequest {
	return domain.TaskSplitRequest{
		ID:         model.ID,
		TaskID:     model.TaskID,
		UserID:     model.UserID,
		Text:       model.Text,
		TemplateID: model.TemplateID,
		Status:     model.Status,
		Result:     model.Result,
		Error:      model.Error,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}
}

