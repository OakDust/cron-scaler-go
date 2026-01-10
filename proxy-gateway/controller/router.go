package controller

import "net/http"

type Router struct {
	controller *Controller
}

func NewRouter(ctrl *Controller) *Router {
	return &Router{controller: ctrl}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := req.Method

	switch {
	case path == "/health" && method == "GET":
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok"}`))

	case path == "/v1/schedules" && method == "POST":
		r.controller.CreateSchedule(w, req)

	case path == "/v1/schedules" && method == "GET":
		r.controller.ListSchedules(w, req)

	case isScheduleWithID(path) && method == "GET":
		r.controller.GetSchedule(w, req)

	case isScheduleWithID(path) && method == "PUT":
		r.controller.UpdateSchedule(w, req)

	case isScheduleWithID(path) && method == "DELETE":
		r.controller.DeleteSchedule(w, req)

	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func isScheduleWithID(path string) bool {
	// Проверяем что путь начинается с /v1/schedules/ и после есть что-то
	if len(path) <= len("/v1/schedules/") {
		return false
	}
	return path[:len("/v1/schedules/")] == "/v1/schedules/"
}
