package controller

import (
	"context"

	scalehandlerv1 "scale-handler/pkg/api/proto/scale-handler"
)

func (c *Controller) Delete(ctx context.Context, req *scalehandlerv1.DeleteRequest) (*scalehandlerv1.DeleteResponse, error) {
	c.logger.Info("Handling Delete request", "id", req.Id)

	// Удаляем расписание через usecase
	err := c.scheduleUC.DeleteSchedule(ctx, req.Id)
	if err != nil {
		c.logger.Error("Failed to delete schedule", "id", req.Id, "error", err)
		return nil, err
	}

	return &scalehandlerv1.DeleteResponse{
		Success: true,
	}, nil
}
