package usecase

import (
	"context"
	"log/slog"

	"scale-handler/internal/domain"
	"scale-handler/internal/repository"
)

type ScheduleUseCase struct {
	repo   repository.ScheduleRepository
	logger *slog.Logger
}

func NewScheduleUseCase(repo repository.ScheduleRepository, logger *slog.Logger) *ScheduleUseCase {
	return &ScheduleUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *ScheduleUseCase) CreateSchedule(ctx context.Context, rules domain.ScheduleRules, application *domain.Application) (*domain.Schedule, error) {
	uc.logger.Debug("Creating schedule", "rules", rules)
	return uc.repo.Create(ctx, rules, application)
}

func (uc *ScheduleUseCase) GetSchedule(ctx context.Context, id string) (*domain.Schedule, error) {
	uc.logger.Debug("Getting schedule", "id", id)
	return uc.repo.GetByID(ctx, id)
}

func (uc *ScheduleUseCase) ListSchedules(ctx context.Context) ([]*domain.Schedule, error) {
	uc.logger.Debug("Listing schedules")
	return uc.repo.List(ctx)
}

func (uc *ScheduleUseCase) UpdateSchedule(ctx context.Context, id string, rules domain.ScheduleRules, application *domain.Application) (*domain.Schedule, error) {
	uc.logger.Debug("Updating schedule", "id", id, "rules", rules)
	return uc.repo.Update(ctx, id, rules, application)
}

func (uc *ScheduleUseCase) DeleteSchedule(ctx context.Context, id string) error {
	uc.logger.Debug("Deleting schedule", "id", id)
	return uc.repo.Delete(ctx, id)
}
