package controller

import (
	"context"

	"scale-handler/internal/controller/converter"
	scalehandlerv1 "scale-handler/pkg/api/proto/scale-handler"
)

func (c *Controller) List(ctx context.Context, req *scalehandlerv1.ListRequest) (*scalehandlerv1.ListResponse, error) {
	c.logger.Info("Handling List request")

	// Получаем список расписаний через usecase
	schedules, err := c.scheduleUC.ListSchedules(ctx)
	if err != nil {
		c.logger.Error("Failed to list schedules", "error", err)
		return nil, err
	}

	// Конвертируем доменные модели в proto
	var protoSchedules []*scalehandlerv1.Schedule
	for _, schedule := range schedules {
		protoSchedules = append(protoSchedules, converter.DomainToProto(schedule))
	}

	return &scalehandlerv1.ListResponse{
		Schedules: protoSchedules,
	}, nil
}
