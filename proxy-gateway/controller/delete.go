package controller

import (
	"net/http"

	scalehandlerv1 "proxy-gateway/pkg/api/proto/scale-handler"

	"github.com/google/uuid"
)

func (c *Controller) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c.logger.Info("Handling delete schedule request")

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
	req := &scalehandlerv1.DeleteRequest{Id: id}
	resp, err := c.grpcClient.Delete(ctx, req)
	if err != nil {
		c.logger.Error("gRPC call failed", "error", err, "id", id)
		writeError(w, http.StatusInternalServerError, "Failed to delete schedule")
		return
	}

	// Возвращаем ответ
	writeJSON(w, http.StatusOK, map[string]bool{
		"success": resp.Success,
	})
}
