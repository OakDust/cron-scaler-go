package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	scalehandlerv1 "proxy-gateway/pkg/api/proto/scale-handler"
)

type CreateScheduleRequest struct {
	Schedule *scalehandlerv1.Schedule `json:"schedule"`
}

func (c *Controller) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c.logger.Info("Handling create schedule request")

	// Проверяем Content-Type
	contentType := r.Header.Get("Content-Type")

	var scheduleReq CreateScheduleRequest
	var err error

	if strings.Contains(contentType, "multipart/form-data") {
		scheduleReq, err = c.parseMultipartForm(r)
	} else {
		scheduleReq, err = c.parseJSONBody(r)
	}

	if err != nil {
		c.logger.Error("Failed to parse request", "error", err)
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err))
		return
	}

	if scheduleReq.Schedule == nil {
		writeError(w, http.StatusBadRequest, "Schedule is required")
		return
	}

	// Простая валидация времени
	if err := c.validateSchedule(scheduleReq.Schedule); err != nil {
		c.logger.Error("Schedule validation failed", "error", err)
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Validation failed: %v", err))
		return
	}

	// Отправляем запрос в gRPC сервис
	req := &scalehandlerv1.CreateRequest{
		Schedule: scheduleReq.Schedule,
	}

	resp, err := c.grpcClient.Create(ctx, req)
	if err != nil {
		c.logger.Error("gRPC call failed", "error", err)
		writeError(w, http.StatusInternalServerError, "Failed to create schedule")
		return
	}

	// Возвращаем ответ
	writeJSON(w, http.StatusCreated, map[string]string{
		"id": resp.Id,
	})
}

func (c *Controller) parseJSONBody(r *http.Request) (CreateScheduleRequest, error) {
	var req CreateScheduleRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return req, fmt.Errorf("failed to read body: %w", err)
	}

	if err := json.Unmarshal(body, &req); err != nil {
		return req, fmt.Errorf("invalid JSON: %w", err)
	}

	return req, nil
}

func (c *Controller) parseMultipartForm(r *http.Request) (CreateScheduleRequest, error) {
	// Парсим multipart форму (максимум 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return CreateScheduleRequest{}, fmt.Errorf("failed to parse form: %w", err)
	}

	// Пытаемся получить JSON из поля "schedule"
	jsonStr := r.FormValue("schedule")
	if jsonStr == "" {
		return CreateScheduleRequest{}, fmt.Errorf("field 'schedule' is required")
	}

	var req CreateScheduleRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		return req, fmt.Errorf("invalid JSON in form field: %w", err)
	}

	return req, nil
}

func (c *Controller) validateSchedule(schedule *scalehandlerv1.Schedule) error {
	// Проверяем формат времени HH:MM
	timeRegex := regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]$`)

	// Проверяем weekdays
	for _, daySchedule := range schedule.Weekdays {
		if daySchedule == nil {
			continue
		}
		for _, tr := range daySchedule.TimeRanges {
			if tr == nil {
				continue
			}
			if !timeRegex.MatchString(tr.From) {
				return fmt.Errorf("invalid time format for 'from': %s", tr.From)
			}
			if !timeRegex.MatchString(tr.To) {
				return fmt.Errorf("invalid time format for 'to': %s", tr.To)
			}

			// Проверяем что from < to
			fromTime, _ := time.Parse("15:04", tr.From)
			toTime, _ := time.Parse("15:04", tr.To)
			if !fromTime.Before(toTime) {
				return fmt.Errorf("'from' time must be before 'to' time: %s - %s", tr.From, tr.To)
			}
		}
	}

	// Проверяем dates (формат DD-MM-YYYY)
	dateRegex := regexp.MustCompile(`^\d{2}-\d{2}-\d{4}$`)
	for date := range schedule.Dates {
		if !dateRegex.MatchString(date) {
			return fmt.Errorf("invalid date format: %s, expected DD-MM-YYYY", date)
		}
	}

	// Проверяем exceptions
	for _, date := range schedule.Exceptions {
		if !dateRegex.MatchString(date) {
			return fmt.Errorf("invalid exception date format: %s, expected DD-MM-YYYY", date)
		}
	}

	return nil
}
