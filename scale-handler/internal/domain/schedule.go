package domain

import (
	"time"
)

type Schedule struct {
	ID        string
	Rules     ScheduleRules
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ScheduleRules struct {
	Weekdays   map[string][]TimeRange `json:"weekdays"`
	Dates      map[string][]TimeRange `json:"dates"`
	Exceptions []string               `json:"exceptions"`
}

type TimeRange struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Replicas int32  `json:"replicas"`
}
