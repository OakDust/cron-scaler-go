package controller

import (
	"net/http"

	scalehandlerv1 "proxy-gateway/pkg/api/proto/scale-handler"
	"proxy-gateway/pkg/schedule"
)

func (c *Controller) ListSchedules(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c.logger.Info("Handling list schedules request")

	// Вызываем gRPC метод
	req := &scalehandlerv1.ListRequest{}
	resp, err := c.grpcClient.List(ctx, req)
	if err != nil {
		c.logger.Error("gRPC call failed", "error", err)
		writeError(w, http.StatusInternalServerError, "Failed to list schedules")
		return
	}

	// Конвертируем в удобный REST-формат
	schedulesDTO := make([]*schedule.ScheduleDTO, len(resp.Schedules))
	for i, s := range resp.Schedules {
		schedulesDTO[i] = schedule.ProtoToDTO(s)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"schedules": schedulesDTO,
	})
}
