package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	scalehandlerv1 "proxy-gateway/pkg/api/proto/scale-handler"
	"proxy-gateway/pkg/schedule"

	"github.com/google/uuid"
)

type UpdateScheduleRequest struct {
	Schedule *schedule.ScheduleDTO `json:"schedule"`
}

func (c *Controller) UpdateSchedule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c.logger.Info("Handling update schedule request")

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

	// Парсим тело запроса
	var req UpdateScheduleRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.logger.Error("Failed to read body", "error", err)
		writeError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		c.logger.Error("Invalid JSON", "error", err)
		writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if req.Schedule == nil {
		writeError(w, http.StatusBadRequest, "Schedule is required")
		return
	}

	// Валидируем schedule
	if err := c.validateScheduleDTO(req.Schedule); err != nil {
		c.logger.Error("Schedule validation failed", "error", err)
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Validation failed: %v", err))
		return
	}

	// Конвертируем в proto и вызываем gRPC
	protoSchedule := schedule.DTOToProto(req.Schedule)
	grpcReq := &scalehandlerv1.UpdateRequest{
		Id:       id,
		Schedule: protoSchedule,
	}

	resp, err := c.grpcClient.Update(ctx, grpcReq)
	if err != nil {
		c.logger.Error("gRPC call failed", "error", err, "id", id)
		writeError(w, http.StatusInternalServerError, "Failed to update schedule")
		return
	}

	// Возвращаем ответ
	writeJSON(w, http.StatusOK, map[string]bool{
		"success": resp.Success,
	})
}
