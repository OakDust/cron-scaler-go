package app

import (
	"fmt"
	"log/slog"
	"net"

	"scale-handler/internal/controller"
	scalehandlerv1 "scale-handler/pkg/api/proto/scale-handler"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	server   *grpc.Server
	listener net.Listener
	logger   *slog.Logger
}

func NewGRPCServer(port string, ctrl *controller.Controller, logger *slog.Logger) (*GRPCServer, error) {
	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()

	// Регистрируем наш сервис
	scalehandlerv1.RegisterScaleHandlerServiceServer(grpcServer, ctrl)

	// Регистрируем gRPC reflection для интроспекции (grpcurl, Postman и т.д.)
	reflection.Register(grpcServer)

	// Создаем listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	return &GRPCServer{
		server:   grpcServer,
		listener: listener,
		logger:   logger,
	}, nil
}

func (s *GRPCServer) Start() error {
	s.logger.Info("Starting gRPC server", "address", s.listener.Addr())
	return s.server.Serve(s.listener)
}

func (s *GRPCServer) Stop() {
	s.logger.Info("Stopping gRPC server")
	s.server.GracefulStop()
}
