package controller

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	scalehandlerv1 "proxy-gateway/pkg/api/proto/scale-handler"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Controller struct {
	logger         *slog.Logger
	grpcClient     scalehandlerv1.ScaleHandlerServiceClient
	grpcConnection *grpc.ClientConn
}

func NewController(logger *slog.Logger, grpcAddr string) (*Controller, error) {
	// Подключаемся к gRPC серверу
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	client := scalehandlerv1.NewScaleHandlerServiceClient(conn)

	return &Controller{
		logger:         logger,
		grpcClient:     client,
		grpcConnection: conn,
	}, nil
}

func (c *Controller) Close() error {
	if c.grpcConnection != nil {
		return c.grpcConnection.Close()
	}
	return nil
}

func (c *Controller) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Schedule endpoints - используем старый синтаксис
	mux.HandleFunc("/v1/schedules", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			c.CreateSchedule(w, r)
		case "GET":
			c.ListSchedules(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	// Для GET/PUT/DELETE с ID
	mux.HandleFunc("/v1/schedules/", func(w http.ResponseWriter, r *http.Request) {
		// Проверяем что после /v1/schedules/ есть что-то
		path := r.URL.Path
		if path == "/v1/schedules/" {
			writeError(w, http.StatusBadRequest, "Schedule ID is required")
			return
		}

		switch r.Method {
		case "GET":
			c.GetSchedule(w, r)
		case "PUT":
			c.UpdateSchedule(w, r)
		case "DELETE":
			c.DeleteSchedule(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok"}`))
	})

	return mux
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
