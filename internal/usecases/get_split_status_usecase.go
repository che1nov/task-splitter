package usecases

import (
	"context"
	"task-splitter/internal/domain"
	"task-splitter/internal/dto"
	"task-splitter/pkg/logger"
)

// GetSplitStatusUseCase Use Case для получения статуса разбивки
type GetSplitStatusUseCase struct {
	splitRequestRepo TaskSplitRequestRepository
	cache           CacheAdapter
	log             logger.Logger
}

// NewGetSplitStatusUseCase создает новый UseCase для получения статуса разбивки
func NewGetSplitStatusUseCase(
	splitRequestRepo TaskSplitRequestRepository,
	cache CacheAdapter,
	log logger.Logger,
) *GetSplitStatusUseCase {
	return &GetSplitStatusUseCase{
		splitRequestRepo: splitRequestRepo,
		cache:           cache,
		log:             log,
	}
}

// Execute получает статус разбивки задачи
func (uc *GetSplitStatusUseCase) Execute(ctx context.Context, requestID string) (dto.SplitStatusOutput, error) {
	uc.log.DebugContext(ctx, "Получение статуса разбивки", "request_id", requestID)

	// Сначала проверяем кэш (если доступен)
	if uc.cache != nil {
		cachedResult, err := uc.cache.Get(ctx, "split_result:"+requestID)
		if err == nil && cachedResult != "" {
			uc.log.DebugContext(ctx, "Результат найден в кэше", "request_id", requestID)
			return dto.SplitStatusOutput{
				RequestID: requestID,
				Status:    domain.SplitStatusCompleted,
				Result:    cachedResult,
			}, nil
		}
	}

	// Получаем из БД
	splitRequest, err := uc.splitRequestRepo.GetSplitRequestByID(ctx, requestID)
	if err != nil {
		uc.log.ErrorContext(ctx, "Запрос на разбивку не найден", "request_id", requestID, "error", err)
		return dto.SplitStatusOutput{}, domain.ErrSplitRequestNotFound
	}

	uc.log.DebugContext(ctx, "Статус разбивки получен", "request_id", requestID, "status", splitRequest.Status)

	return dto.SplitStatusOutput{
		RequestID: splitRequest.ID,
		Status:    splitRequest.Status,
		Result:    splitRequest.Result,
		Error:     splitRequest.Error,
		CreatedAt: splitRequest.CreatedAt,
		UpdatedAt: splitRequest.UpdatedAt,
	}, nil
}

