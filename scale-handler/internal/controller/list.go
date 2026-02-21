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

	items := make([]*scalehandlerv1.ScheduleWithApplication, len(schedules))
	for i, s := range schedules {
		items[i] = &scalehandlerv1.ScheduleWithApplication{
			Schedule:    converter.DomainToProto(s),
			Application: converter.ApplicationToProto(s.Application),
		}
	}

	return &scalehandlerv1.ListResponse{
		Items: items,
	}, nil
}
