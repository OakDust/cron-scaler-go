package controller

import (
	"context"

	"scale-handler/internal/controller/converter"
	scalehandlerv1 "scale-handler/pkg/api/proto/scale-handler"
)

func (c *Controller) Update(ctx context.Context, req *scalehandlerv1.UpdateRequest) (*scalehandlerv1.UpdateResponse, error) {
	c.logger.Info("Handling Update request", "id", req.Id)

	// Конвертируем proto в доменные правила
	rules := converter.ProtoToDomainRules(req.Schedule)

	// Обновляем расписание через usecase
	_, err := c.scheduleUC.UpdateSchedule(ctx, req.Id, rules)
	if err != nil {
		c.logger.Error("Failed to update schedule", "id", req.Id, "error", err)
		return nil, err
	}

	return &scalehandlerv1.UpdateResponse{
		Success: true,
	}, nil
}
