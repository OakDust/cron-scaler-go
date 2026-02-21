package controller

import (
	"context"

	"scale-handler/internal/controller/converter"
	scalehandlerv1 "scale-handler/pkg/api/proto/scale-handler"
)

func (c *Controller) Create(ctx context.Context, req *scalehandlerv1.CreateRequest) (*scalehandlerv1.CreateResponse, error) {
	c.logger.Info("Handling Create request")

	rules := converter.ProtoToDomainRules(req.Schedule)
	application := converter.ProtoToApplication(req.Application)

	schedule, err := c.scheduleUC.CreateSchedule(ctx, rules, application)
	if err != nil {
		c.logger.Error("Failed to create schedule", "error", err)
		return nil, err
	}

	if c.k8sReconciler != nil {
		if err := c.k8sReconciler.CreateResources(ctx, schedule); err != nil {
			c.logger.Error("Failed to create K8s resources", "id", schedule.ID, "error", err)
		}
	}

	return &scalehandlerv1.CreateResponse{
		Id: schedule.ID,
	}, nil
}
