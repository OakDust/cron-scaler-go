package controller

import (
	"context"

	"scale-handler/internal/controller/converter"
	scalehandlerv1 "scale-handler/pkg/api/proto/scale-handler"
)

func (c *Controller) Create(ctx context.Context, req *scalehandlerv1.CreateRequest) (*scalehandlerv1.CreateResponse, error) {
	c.logger.Info("Handling Create request")

	// Конвертируем proto в доменные правила
	rules := converter.ProtoToDomainRules(req.Schedule)

	// Создаем расписание через usecase
	schedule, err := c.scheduleUC.CreateSchedule(ctx, rules)
	if err != nil {
		c.logger.Error("Failed to create schedule", "error", err)
		return nil, err
	}

	return &scalehandlerv1.CreateResponse{
		Id: schedule.ID,
	}, nil
}
