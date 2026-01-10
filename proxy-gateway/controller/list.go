package controller

import (
	"net/http"

	scalehandlerv1 "proxy-gateway/pkg/api/proto/scale-handler"
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

	// Возвращаем ответ
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"schedules": resp.Schedules,
	})
}
