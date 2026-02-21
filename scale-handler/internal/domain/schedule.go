package domain

import (
	"time"
)

type Schedule struct {
	ID          string
	Rules       ScheduleRules
	Application *Application
	CreatedAt   time.Time
	UpdatedAt   time.Time
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

type Application struct {
	Containers []Container `json:"containers"`
}

type Container struct {
	Name            string           `json:"name"`
	Image           string           `json:"image"`
	Ports           []ContainerPort  `json:"ports,omitempty"`
	Env             []EnvVar         `json:"env,omitempty"`
	Resources       *Resources       `json:"resources,omitempty"`
	LivenessProbe   *Probe           `json:"livenessProbe,omitempty"`
	ReadinessProbe  *Probe           `json:"readinessProbe,omitempty"`
}

type ContainerPort struct {
	ContainerPort int    `json:"containerPort"`
	Protocol     string `json:"protocol,omitempty"`
}

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Resources struct {
	Requests *ResourceQuantity `json:"requests,omitempty"`
	Limits   *ResourceQuantity `json:"limits,omitempty"`
}

type ResourceQuantity struct {
	Memory string `json:"memory,omitempty"`
	CPU    string `json:"cpu,omitempty"`
}

type Probe struct {
	HTTPGet             *HTTPGetAction `json:"httpGet,omitempty"`
	InitialDelaySeconds int32          `json:"initialDelaySeconds,omitempty"`
	PeriodSeconds       int32          `json:"periodSeconds,omitempty"`
}

type HTTPGetAction struct {
	Path string `json:"path"`
	Port int32  `json:"port"`
}
