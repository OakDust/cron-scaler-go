package controller

import (
	"log/slog"

	"scale-handler/internal/usecase"
	scalehandlerv1 "scale-handler/pkg/api/proto/scale-handler"
)

// Убедимся что Controller реализует интерфейс
var _ scalehandlerv1.ScaleHandlerServiceServer = (*Controller)(nil)

type Controller struct {
	scalehandlerv1.UnimplementedScaleHandlerServiceServer
	scheduleUC *usecase.ScheduleUseCase
	logger     *slog.Logger
}

func NewController(scheduleUC *usecase.ScheduleUseCase, logger *slog.Logger) *Controller {
	return &Controller{
		scheduleUC: scheduleUC,
		logger:     logger,
	}
}
