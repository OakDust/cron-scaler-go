package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	scalehandlerv1 "proxy-gateway/pkg/api/proto/scale-handler"
	"proxy-gateway/pkg/schedule"
)

type CreateScheduleRequest struct {
	Schedule    *schedule.ScheduleDTO    `json:"schedule"`
	Application *schedule.ApplicationDTO `json:"application"`
}

// CreateSchedule godoc
// @Summary      Создать расписание
// @Description  Создаёт новое расписание масштабирования с указанием weekdays, dates, exceptions и application
// @Tags         schedules
// @Accept       json
// @Produce      json
// @Param        body  body  CreateScheduleRequest  true  "Schedule and Application"
// @Success      201   {object}  map[string]string  "id"
// @Failure      400  {object}  map[string]string  "error"
// @Failure      500  {object}  map[string]string  "error"
// @Router       /v1/schedules [post]
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

	// Валидация формата
	if err := c.validateScheduleDTO(scheduleReq.Schedule); err != nil {
		c.logger.Error("Schedule validation failed", "error", err)
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Validation failed: %v", err))
		return
	}

	protoSchedule := schedule.DTOToProto(scheduleReq.Schedule)
	protoApp := schedule.ApplicationDTOToProto(scheduleReq.Application)
	req := &scalehandlerv1.CreateRequest{
		Schedule:    protoSchedule,
		Application: protoApp,
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

	jsonStr := r.FormValue("schedule")
	if jsonStr == "" {
		jsonStr = r.FormValue("data")
	}
	if jsonStr == "" {
		return CreateScheduleRequest{}, fmt.Errorf("field 'schedule' or 'data' is required")
	}

	var req CreateScheduleRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		return req, fmt.Errorf("invalid JSON in form field: %w", err)
	}

	return req, nil
}

func (c *Controller) validateScheduleDTO(s *schedule.ScheduleDTO) error {
	// Проверяем формат времени HH:MM
	timeRegex := regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]$`)
	// Формат даты ISO 8601: YYYY-MM-DD
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

	// Проверяем weekdays + пересечение интервалов
	for day, ranges := range s.Weekdays {
		if len(ranges) == 0 {
			continue
		}
		// базовая проверка формата + from < to
		for _, tr := range ranges {
			if !timeRegex.MatchString(tr.From) {
				return fmt.Errorf("invalid time format for 'from' in %s: %s", day, tr.From)
			}
			if !timeRegex.MatchString(tr.To) {
				return fmt.Errorf("invalid time format for 'to' in %s: %s", day, tr.To)
			}
			fromTime, _ := time.Parse("15:04", tr.From)
			toTime, _ := time.Parse("15:04", tr.To)
			if !fromTime.Before(toTime) {
				return fmt.Errorf("'from' time must be before 'to' time: %s - %s", tr.From, tr.To)
			}
		}
		// запрет пересечения интервалов
		if err := ensureNoOverlapDTO(day, ranges); err != nil {
			return err
		}
	}

	// Проверяем dates (формат YYYY-MM-DD) и пересечение интервалов для каждой даты
	for date, ranges := range s.Dates {
		if !dateRegex.MatchString(date) {
			return fmt.Errorf("invalid date format: %s, expected YYYY-MM-DD", date)
		}
		if len(ranges) == 0 {
			continue
		}
		for _, tr := range ranges {
			if !timeRegex.MatchString(tr.From) {
				return fmt.Errorf("invalid time format for 'from' in date %s: %s", date, tr.From)
			}
			if !timeRegex.MatchString(tr.To) {
				return fmt.Errorf("invalid time format for 'to' in date %s: %s", date, tr.To)
			}
			fromTime, _ := time.Parse("15:04", tr.From)
			toTime, _ := time.Parse("15:04", tr.To)
			if !fromTime.Before(toTime) {
				return fmt.Errorf("'from' time must be before 'to' time in date %s: %s - %s", date, tr.From, tr.To)
			}
		}
		if err := ensureNoOverlapDTO(date, ranges); err != nil {
			return err
		}
	}

	// Проверяем exceptions
	for _, date := range s.Exceptions {
		if !dateRegex.MatchString(date) {
			return fmt.Errorf("invalid exception date format: %s, expected YYYY-MM-DD", date)
		}
	}

	return nil
}

func (c *Controller) validateSchedule(s *scalehandlerv1.Schedule) error {
	timeRegex := regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]$`)
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

	for name, daySchedule := range s.Weekdays {
		if daySchedule == nil {
			continue
		}
		if len(daySchedule.TimeRanges) == 0 {
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
			fromTime, _ := time.Parse("15:04", tr.From)
			toTime, _ := time.Parse("15:04", tr.To)
			if !fromTime.Before(toTime) {
				return fmt.Errorf("'from' time must be before 'to' time: %s - %s", tr.From, tr.To)
			}
		}
		if err := ensureNoOverlapProto(name, daySchedule.TimeRanges); err != nil {
			return err
		}
	}

	for date, daySchedule := range s.Dates {
		if !dateRegex.MatchString(date) {
			return fmt.Errorf("invalid date format: %s, expected YYYY-MM-DD", date)
		}
		if daySchedule == nil || len(daySchedule.TimeRanges) == 0 {
			continue
		}
		for _, tr := range daySchedule.TimeRanges {
			if tr == nil {
				continue
			}
			if !timeRegex.MatchString(tr.From) {
				return fmt.Errorf("invalid time format for 'from' in date %s: %s", date, tr.From)
			}
			if !timeRegex.MatchString(tr.To) {
				return fmt.Errorf("invalid time format for 'to' in date %s: %s", date, tr.To)
			}
			fromTime, _ := time.Parse("15:04", tr.From)
			toTime, _ := time.Parse("15:04", tr.To)
			if !fromTime.Before(toTime) {
				return fmt.Errorf("'from' time must be before 'to' time in date %s: %s - %s", date, tr.From, tr.To)
			}
		}
		if err := ensureNoOverlapProto(date, daySchedule.TimeRanges); err != nil {
			return err
		}
	}

	for _, date := range s.Exceptions {
		if !dateRegex.MatchString(date) {
			return fmt.Errorf("invalid exception date format: %s, expected YYYY-MM-DD", date)
		}
	}

	return nil
}

func ensureNoOverlapDTO(name string, ranges []schedule.TimeRangeDTO) error {
	if len(ranges) < 2 {
		return nil
	}
	sorted := make([]schedule.TimeRangeDTO, len(ranges))
	copy(sorted, ranges)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].From < sorted[j].From
	})
	prevFrom, _ := time.Parse("15:04", sorted[0].From)
	prevTo, _ := time.Parse("15:04", sorted[0].To)
	for i := 1; i < len(sorted); i++ {
		curFrom, _ := time.Parse("15:04", sorted[i].From)
		curTo, _ := time.Parse("15:04", sorted[i].To)
		if curFrom.Before(prevTo) {
			return fmt.Errorf("time ranges overlap in %s: %s-%s and %s-%s",
				name, prevFrom.Format("15:04"), prevTo.Format("15:04"), curFrom.Format("15:04"), curTo.Format("15:04"))
		}
		prevFrom, prevTo = curFrom, curTo
	}
	return nil
}

func ensureNoOverlapProto(name string, ranges []*scalehandlerv1.TimeRange) error {
	if len(ranges) < 2 {
		return nil
	}
	type trWithParsed struct {
		from time.Time
		to   time.Time
		raw  *scalehandlerv1.TimeRange
	}
	tmp := make([]trWithParsed, 0, len(ranges))
	for _, r := range ranges {
		if r == nil {
			continue
		}
		fromTime, err1 := time.Parse("15:04", r.From)
		toTime, err2 := time.Parse("15:04", r.To)
		if err1 != nil || err2 != nil {
			continue
		}
		tmp = append(tmp, trWithParsed{from: fromTime, to: toTime, raw: r})
	}
	if len(tmp) < 2 {
		return nil
	}
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].from.Before(tmp[j].from)
	})
	prev := tmp[0]
	for i := 1; i < len(tmp); i++ {
		cur := tmp[i]
		if cur.from.Before(prev.to) {
			return fmt.Errorf("time ranges overlap in %s: %s-%s and %s-%s",
				name,
				prev.from.Format("15:04"), prev.to.Format("15:04"),
				cur.from.Format("15:04"), cur.to.Format("15:04"),
			)
		}
		prev = cur
	}
	return nil
}
