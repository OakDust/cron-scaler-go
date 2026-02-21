package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"scale-handler/internal/domain"

	"github.com/jmoiron/sqlx"
)

func (r *ScheduleRepository) CheckConnection(ctx context.Context) error {
	// Проверим что таблица существует
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'schedules'
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check table existence: %w", err)
	}

	if !exists {
		return fmt.Errorf("table 'public.schedules' does not exist")
	}

	return nil
}

type ScheduleRepository struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewScheduleRepository(db *sqlx.DB, logger *slog.Logger) *ScheduleRepository {
	return &ScheduleRepository{
		db:     db,
		logger: logger,
	}
}

func (r *ScheduleRepository) Create(ctx context.Context, rules domain.ScheduleRules, application *domain.Application) (*domain.Schedule, error) {
	query := `
		INSERT INTO public.schedules (rules, application)
		VALUES ($1::jsonb, $2::jsonb)
		RETURNING id, rules, application, created_at, updated_at
	`

	rulesJSON, err := json.Marshal(rules)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal rules: %w", err)
	}

	var appArg interface{}
	if application != nil {
		b, _ := json.Marshal(application)
		appArg = string(b)
	}

	var schedule domain.Schedule
	var rulesBytes, appBytes []byte

	err = r.db.QueryRowContext(ctx, query, string(rulesJSON), appArg).Scan(
		&schedule.ID,
		&rulesBytes,
		&appBytes,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	if err := json.Unmarshal(rulesBytes, &schedule.Rules); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rules: %w", err)
	}
	if len(appBytes) > 0 {
		if err := json.Unmarshal(appBytes, &schedule.Application); err != nil {
			return nil, fmt.Errorf("failed to unmarshal application: %w", err)
		}
	}

	return &schedule, nil
}

func (r *ScheduleRepository) GetByID(ctx context.Context, id string) (*domain.Schedule, error) {
	query := `
		SELECT id, rules, application, created_at, updated_at
		FROM schedules
		WHERE id = $1
	`

	var schedule domain.Schedule
	var rulesBytes, appBytes []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&schedule.ID,
		&rulesBytes,
		&appBytes,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("schedule not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	if err := json.Unmarshal(rulesBytes, &schedule.Rules); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rules: %w", err)
	}
	if len(appBytes) > 0 {
		if err := json.Unmarshal(appBytes, &schedule.Application); err != nil {
			return nil, fmt.Errorf("failed to unmarshal application: %w", err)
		}
	}

	return &schedule, nil
}

func (r *ScheduleRepository) List(ctx context.Context) ([]*domain.Schedule, error) {
	query := `
		SELECT id, rules, application, created_at, updated_at
		FROM schedules
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list schedules: %w", err)
	}
	defer rows.Close()

	var schedules []*domain.Schedule

	for rows.Next() {
		var schedule domain.Schedule
		var rulesBytes, appBytes []byte

		if err := rows.Scan(
			&schedule.ID,
			&rulesBytes,
			&appBytes,
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %w", err)
		}

		if err := json.Unmarshal(rulesBytes, &schedule.Rules); err != nil {
			return nil, fmt.Errorf("failed to unmarshal rules: %w", err)
		}
		if len(appBytes) > 0 {
			if err := json.Unmarshal(appBytes, &schedule.Application); err != nil {
				return nil, fmt.Errorf("failed to unmarshal application: %w", err)
			}
		}

		schedules = append(schedules, &schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return schedules, nil
}

func (r *ScheduleRepository) Update(ctx context.Context, id string, rules domain.ScheduleRules, application *domain.Application) (*domain.Schedule, error) {
	query := `
		UPDATE schedules
		SET rules = $1, application = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING id, rules, application, created_at, updated_at
	`

	rulesJSON, err := json.Marshal(rules)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal rules: %w", err)
	}

	var appArg interface{}
	if application != nil {
		b, _ := json.Marshal(application)
		appArg = string(b)
	}

	var schedule domain.Schedule
	var rulesBytes, appBytes []byte

	err = r.db.QueryRowContext(ctx, query, rulesJSON, appArg, id).Scan(
		&schedule.ID,
		&rulesBytes,
		&appBytes,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("schedule not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to update schedule: %w", err)
	}

	if err := json.Unmarshal(rulesBytes, &schedule.Rules); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rules: %w", err)
	}
	if len(appBytes) > 0 {
		if err := json.Unmarshal(appBytes, &schedule.Application); err != nil {
			return nil, fmt.Errorf("failed to unmarshal application: %w", err)
		}
	}

	return &schedule, nil
}

func (r *ScheduleRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM schedules
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete schedule: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("schedule not found: %w", domain.ErrNotFound)
	}

	return nil
}
