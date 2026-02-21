package controller

import (
	"net/http"

	scalehandlerv1 "proxy-gateway/pkg/api/proto/scale-handler"
	"proxy-gateway/pkg/schedule"
)

// ListSchedules godoc
// @Summary      Список расписаний
// @Description  Возвращает все расписания
// @Tags         schedules
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "items"
// @Failure      500  {object}  map[string]string  "error"
// @Router       /v1/schedules [get]
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

	items := make([]map[string]interface{}, len(resp.Items))
	for i, item := range resp.Items {
		scheduleDTO := schedule.ProtoToDTO(item.Schedule)
		appDTO := schedule.ProtoToApplicationDTO(item.Application)
		items[i] = map[string]interface{}{
			"schedule":    scheduleDTO,
			"application": appDTO,
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"items": items,
	})
}
