package controller

import (
	"log/slog"

	"scale-handler/internal/k8s"
	"scale-handler/internal/usecase"
	scalehandlerv1 "scale-handler/pkg/api/proto/scale-handler"
)

var _ scalehandlerv1.ScaleHandlerServiceServer = (*Controller)(nil)

type Controller struct {
	scalehandlerv1.UnimplementedScaleHandlerServiceServer
	scheduleUC    *usecase.ScheduleUseCase
	k8sReconciler *k8s.Reconciler
	logger        *slog.Logger
}

func NewController(scheduleUC *usecase.ScheduleUseCase, k8sReconciler *k8s.Reconciler, logger *slog.Logger) *Controller {
	return &Controller{
		scheduleUC:    scheduleUC,
		k8sReconciler: k8sReconciler,
		logger:        logger,
	}
}
