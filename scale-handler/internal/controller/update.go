package controller

import (
	"context"

	"scale-handler/internal/controller/converter"
	scalehandlerv1 "scale-handler/pkg/api/proto/scale-handler"
)

func (c *Controller) Update(ctx context.Context, req *scalehandlerv1.UpdateRequest) (*scalehandlerv1.UpdateResponse, error) {
	c.logger.Info("Handling Update request", "id", req.Id)

	rules := converter.ProtoToDomainRules(req.Schedule)
	application := converter.ProtoToApplication(req.Application)

	schedule, err := c.scheduleUC.UpdateSchedule(ctx, req.Id, rules, application)
	if err != nil {
		c.logger.Error("Failed to update schedule", "id", req.Id, "error", err)
		return nil, err
	}

	if c.k8sReconciler != nil {
		if err := c.k8sReconciler.UpdateResources(ctx, schedule); err != nil {
			c.logger.Error("Failed to update K8s resources", "id", schedule.ID, "error", err)
		}
	}

	return &scalehandlerv1.UpdateResponse{
		Success: true,
	}, nil
}
