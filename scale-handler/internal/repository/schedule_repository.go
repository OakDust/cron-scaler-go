package repository

import (
	"context"

	"scale-handler/internal/domain"
)

type ScheduleRepository interface {
	Create(ctx context.Context, rules domain.ScheduleRules, application *domain.Application) (*domain.Schedule, error)
	GetByID(ctx context.Context, id string) (*domain.Schedule, error)
	List(ctx context.Context) ([]*domain.Schedule, error)
	Update(ctx context.Context, id string, rules domain.ScheduleRules, application *domain.Application) (*domain.Schedule, error)
	Delete(ctx context.Context, id string) error
}
