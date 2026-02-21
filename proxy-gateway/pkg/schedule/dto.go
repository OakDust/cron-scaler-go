package schedule

// CreateScheduleRequestDTO - REST API формат (example-schedule.json)
type CreateScheduleRequestDTO struct {
	Schedule    *ScheduleDTO    `json:"schedule"`
	Application *ApplicationDTO `json:"application"`
}

// ScheduleDTO - расписание масштабирования
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

// ApplicationDTO - контейнеры для Deployment
type ApplicationDTO struct {
	Containers []ContainerDTO `json:"containers"`
}

type ContainerDTO struct {
	Name      string            `json:"name"`
	Image     string            `json:"image"`
	Ports     []ContainerPortDTO `json:"ports,omitempty"`
	Env       []EnvVarDTO       `json:"env,omitempty"`
	Resources *ResourcesDTO     `json:"resources,omitempty"`
	LivenessProbe  *ProbeDTO `json:"livenessProbe,omitempty"`
	ReadinessProbe *ProbeDTO `json:"readinessProbe,omitempty"`
}

type ContainerPortDTO struct {
	ContainerPort int32  `json:"containerPort"`
	Protocol     string `json:"protocol,omitempty"` // TCP, UDP
}

type EnvVarDTO struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ResourcesDTO struct {
	Requests *ResourceQuantityDTO `json:"requests,omitempty"`
	Limits   *ResourceQuantityDTO `json:"limits,omitempty"`
}

type ResourceQuantityDTO struct {
	Memory string `json:"memory,omitempty"`
	CPU    string `json:"cpu,omitempty"`
}

type ProbeDTO struct {
	HTTPGet             *HTTPGetActionDTO `json:"httpGet,omitempty"`
	InitialDelaySeconds  int32             `json:"initialDelaySeconds,omitempty"`
	PeriodSeconds        int32             `json:"periodSeconds,omitempty"`
}

type HTTPGetActionDTO struct {
	Path string `json:"path"`
	Port int32  `json:"port"`
}
