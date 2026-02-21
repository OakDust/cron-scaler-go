package schedule

// ScheduleDTO - REST API формат (example-schedule.json)
// Удобный формат: нативные JSON объекты, прямые массивы time ranges
type ScheduleDTO struct {
	Weekdays   map[string][]TimeRangeDTO `json:"weekdays"`
	Dates      map[string][]TimeRangeDTO `json:"dates"`
	Exceptions []string                  `json:"exceptions"`
}

type TimeRangeDTO struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Replicas int32  `json:"replicas"`
}
