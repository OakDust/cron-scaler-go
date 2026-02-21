package controller

import (
	"context"

	"scale-handler/internal/controller/converter"
	scalehandlerv1 "scale-handler/pkg/api/proto/scale-handler"
)

func (c *Controller) Get(ctx context.Context, req *scalehandlerv1.GetRequest) (*scalehandlerv1.GetResponse, error) {
	c.logger.Info("Handling Get request", "id", req.Id)

	// Получаем расписание через usecase
	schedule, err := c.scheduleUC.GetSchedule(ctx, req.Id)
	if err != nil {
		c.logger.Error("Failed to get schedule", "id", req.Id, "error", err)
		return nil, err
	}

	protoSchedule := converter.DomainToProto(schedule)
	protoApplication := converter.ApplicationToProto(schedule.Application)

	return &scalehandlerv1.GetResponse{
		Schedule:    protoSchedule,
		Application: protoApplication,
	}, nil
}
