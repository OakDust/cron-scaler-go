package schedule

import (
	scalehandlerv1 "proxy-gateway/pkg/api/proto/scale-handler"
)

// DTOToProto конвертирует REST DTO в gRPC protobuf
func DTOToProto(dto *ScheduleDTO) *scalehandlerv1.Schedule {
	if dto == nil {
		return nil
	}

	proto := &scalehandlerv1.Schedule{
		Weekdays:   make(map[string]*scalehandlerv1.Schedule_DaySchedule),
		Dates:      make(map[string]*scalehandlerv1.Schedule_DaySchedule),
		Exceptions: dto.Exceptions,
	}

	for day, ranges := range dto.Weekdays {
		proto.Weekdays[day] = &scalehandlerv1.Schedule_DaySchedule{
			TimeRanges: timeRangesToProto(ranges),
		}
	}

	for date, ranges := range dto.Dates {
		proto.Dates[date] = &scalehandlerv1.Schedule_DaySchedule{
			TimeRanges: timeRangesToProto(ranges),
		}
	}

	return proto
}

// ProtoToDTO конвертирует gRPC protobuf в REST DTO
func ProtoToDTO(proto *scalehandlerv1.Schedule) *ScheduleDTO {
	if proto == nil {
		return nil
	}

	dto := &ScheduleDTO{
		Weekdays:   make(map[string][]TimeRangeDTO),
		Dates:      make(map[string][]TimeRangeDTO),
		Exceptions: proto.Exceptions,
	}

	for day, daySchedule := range proto.Weekdays {
		if daySchedule != nil {
			dto.Weekdays[day] = timeRangesToDTO(daySchedule.TimeRanges)
		} else {
			dto.Weekdays[day] = []TimeRangeDTO{}
		}
	}

	for date, daySchedule := range proto.Dates {
		if daySchedule != nil {
			dto.Dates[date] = timeRangesToDTO(daySchedule.TimeRanges)
		} else {
			dto.Dates[date] = []TimeRangeDTO{}
		}
	}

	return dto
}

func timeRangesToProto(ranges []TimeRangeDTO) []*scalehandlerv1.TimeRange {
	if ranges == nil {
		return nil
	}
	result := make([]*scalehandlerv1.TimeRange, len(ranges))
	for i, r := range ranges {
		result[i] = &scalehandlerv1.TimeRange{
			From:     r.From,
			To:       r.To,
			Replicas: r.Replicas,
		}
	}
	return result
}

func timeRangesToDTO(ranges []*scalehandlerv1.TimeRange) []TimeRangeDTO {
	if ranges == nil {
		return nil
	}
	result := make([]TimeRangeDTO, len(ranges))
	for i, r := range ranges {
		if r != nil {
			result[i] = TimeRangeDTO{From: r.From, To: r.To, Replicas: r.Replicas}
		}
	}
	return result
}

// ApplicationDTOToProto конвертирует Application DTO в proto
func ApplicationDTOToProto(dto *ApplicationDTO) *scalehandlerv1.Application {
	if dto == nil {
		return nil
	}
	proto := &scalehandlerv1.Application{
		Containers: make([]*scalehandlerv1.Container, len(dto.Containers)),
	}
	for i, c := range dto.Containers {
		proto.Containers[i] = containerDTOToProto(&c)
	}
	return proto
}

// ProtoToApplicationDTO конвертирует proto Application в DTO
func ProtoToApplicationDTO(proto *scalehandlerv1.Application) *ApplicationDTO {
	if proto == nil {
		return nil
	}
	dto := &ApplicationDTO{
		Containers: make([]ContainerDTO, len(proto.Containers)),
	}
	for i, c := range proto.Containers {
		if c != nil {
			dto.Containers[i] = *containerProtoToDTO(c)
		}
	}
	return dto
}

func containerDTOToProto(c *ContainerDTO) *scalehandlerv1.Container {
	if c == nil {
		return nil
	}
	proto := &scalehandlerv1.Container{
		Name:  c.Name,
		Image: c.Image,
	}
	if len(c.Ports) > 0 {
		proto.Ports = make([]*scalehandlerv1.ContainerPort, len(c.Ports))
		for i, p := range c.Ports {
			proto.Ports[i] = &scalehandlerv1.ContainerPort{
				ContainerPort: p.ContainerPort,
				Protocol:      p.Protocol,
			}
		}
	}
	if len(c.Env) > 0 {
		proto.Env = make([]*scalehandlerv1.EnvVar, len(c.Env))
		for i, e := range c.Env {
			proto.Env[i] = &scalehandlerv1.EnvVar{Name: e.Name, Value: e.Value}
		}
	}
	if c.Resources != nil {
		proto.Resources = &scalehandlerv1.Resources{}
		if c.Resources.Requests != nil {
			proto.Resources.Requests = &scalehandlerv1.ResourceQuantity{
				Memory: c.Resources.Requests.Memory,
				Cpu:    c.Resources.Requests.CPU,
			}
		}
		if c.Resources.Limits != nil {
			proto.Resources.Limits = &scalehandlerv1.ResourceQuantity{
				Memory: c.Resources.Limits.Memory,
				Cpu:    c.Resources.Limits.CPU,
			}
		}
	}
	if c.LivenessProbe != nil && c.LivenessProbe.HTTPGet != nil {
		proto.LivenessProbe = &scalehandlerv1.Probe{
			HttpGet:             &scalehandlerv1.HttpGetAction{Path: c.LivenessProbe.HTTPGet.Path, Port: c.LivenessProbe.HTTPGet.Port},
			InitialDelaySeconds: c.LivenessProbe.InitialDelaySeconds,
			PeriodSeconds:       c.LivenessProbe.PeriodSeconds,
		}
	}
	if c.ReadinessProbe != nil && c.ReadinessProbe.HTTPGet != nil {
		proto.ReadinessProbe = &scalehandlerv1.Probe{
			HttpGet:             &scalehandlerv1.HttpGetAction{Path: c.ReadinessProbe.HTTPGet.Path, Port: c.ReadinessProbe.HTTPGet.Port},
			InitialDelaySeconds: c.ReadinessProbe.InitialDelaySeconds,
			PeriodSeconds:       c.ReadinessProbe.PeriodSeconds,
		}
	}
	return proto
}

func containerProtoToDTO(proto *scalehandlerv1.Container) *ContainerDTO {
	if proto == nil {
		return nil
	}
	c := &ContainerDTO{
		Name:  proto.Name,
		Image: proto.Image,
	}
	if len(proto.Ports) > 0 {
		c.Ports = make([]ContainerPortDTO, len(proto.Ports))
		for i, p := range proto.Ports {
			if p != nil {
				c.Ports[i] = ContainerPortDTO{ContainerPort: p.ContainerPort, Protocol: p.Protocol}
			}
		}
	}
	if len(proto.Env) > 0 {
		c.Env = make([]EnvVarDTO, len(proto.Env))
		for i, e := range proto.Env {
			if e != nil {
				c.Env[i] = EnvVarDTO{Name: e.Name, Value: e.Value}
			}
		}
	}
	if proto.Resources != nil {
		c.Resources = &ResourcesDTO{}
		if proto.Resources.Requests != nil {
			c.Resources.Requests = &ResourceQuantityDTO{
				Memory: proto.Resources.Requests.Memory,
				CPU:    proto.Resources.Requests.Cpu,
			}
		}
		if proto.Resources.Limits != nil {
			c.Resources.Limits = &ResourceQuantityDTO{
				Memory: proto.Resources.Limits.Memory,
				CPU:    proto.Resources.Limits.Cpu,
			}
		}
	}
	if proto.LivenessProbe != nil && proto.LivenessProbe.HttpGet != nil {
		c.LivenessProbe = &ProbeDTO{
			HTTPGet:             &HTTPGetActionDTO{Path: proto.LivenessProbe.HttpGet.Path, Port: proto.LivenessProbe.HttpGet.Port},
			InitialDelaySeconds: proto.LivenessProbe.InitialDelaySeconds,
			PeriodSeconds:       proto.LivenessProbe.PeriodSeconds,
		}
	}
	if proto.ReadinessProbe != nil && proto.ReadinessProbe.HttpGet != nil {
		c.ReadinessProbe = &ProbeDTO{
			HTTPGet:             &HTTPGetActionDTO{Path: proto.ReadinessProbe.HttpGet.Path, Port: proto.ReadinessProbe.HttpGet.Port},
			InitialDelaySeconds: proto.ReadinessProbe.InitialDelaySeconds,
			PeriodSeconds:       proto.ReadinessProbe.PeriodSeconds,
		}
	}
	return c
}
