package controller

import (
	"net/http"
	"strings"

	scalehandlerv1 "proxy-gateway/pkg/api/proto/scale-handler"
	"proxy-gateway/pkg/schedule"

	"github.com/google/uuid"
)

// GetSchedule godoc
// @Summary      Получить расписание
// @Description  Получает расписание по ID
// @Tags         schedules
// @Produce      json
// @Param        id   path      string  true  "Schedule UUID"
// @Success      200  {object}  map[string]interface{}  "schedule, application"
// @Failure      400  {object}  map[string]string  "error"
// @Failure      404  {object}  map[string]string  "error"
// @Router       /v1/schedules/{id} [get]
func (c *Controller) GetSchedule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c.logger.Info("Handling get schedule request")

	// Извлекаем ID из пути
	id := extractIDFromPath(r.URL.Path)
	if id == "" {
		writeError(w, http.StatusBadRequest, "Schedule ID is required")
		return
	}

	// Проверяем UUID
	if _, err := uuid.Parse(id); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	// Вызываем gRPC метод
	req := &scalehandlerv1.GetRequest{Id: id}
	resp, err := c.grpcClient.Get(ctx, req)
	if err != nil {
		c.logger.Error("gRPC call failed", "error", err, "id", id)
		writeError(w, http.StatusNotFound, "Schedule not found")
		return
	}

	scheduleDTO := schedule.ProtoToDTO(resp.Schedule)
	appDTO := schedule.ProtoToApplicationDTO(resp.Application)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"schedule":    scheduleDTO,
		"application": appDTO,
	})
}

func extractIDFromPath(path string) string {
	// Ожидаем путь вида /v1/schedules/{id}
	parts := strings.Split(path, "/")
	if len(parts) >= 4 {
		return parts[3]
	}
	return ""
}
