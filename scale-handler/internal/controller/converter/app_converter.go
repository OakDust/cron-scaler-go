package converter

import (
	"scale-handler/internal/domain"
	scalehandlerv1 "scale-handler/pkg/api/proto/scale-handler"
)

func ApplicationToProto(app *domain.Application) *scalehandlerv1.Application {
	if app == nil {
		return nil
	}
	proto := &scalehandlerv1.Application{
		Containers: make([]*scalehandlerv1.Container, len(app.Containers)),
	}
	for i, c := range app.Containers {
		proto.Containers[i] = containerToProto(&c)
	}
	return proto
}

func ProtoToApplication(proto *scalehandlerv1.Application) *domain.Application {
	if proto == nil {
		return nil
	}
	app := &domain.Application{
		Containers: make([]domain.Container, len(proto.Containers)),
	}
	for i, c := range proto.Containers {
		if c != nil {
			app.Containers[i] = *containerToDomain(c)
		}
	}
	return app
}

func containerToProto(c *domain.Container) *scalehandlerv1.Container {
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
				ContainerPort: int32(p.ContainerPort),
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
		if proto.Resources.Requests == nil && proto.Resources.Limits == nil {
			proto.Resources = nil
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

func containerToDomain(proto *scalehandlerv1.Container) *domain.Container {
	if proto == nil {
		return nil
	}
	c := &domain.Container{
		Name:  proto.Name,
		Image: proto.Image,
	}
	if len(proto.Ports) > 0 {
		c.Ports = make([]domain.ContainerPort, len(proto.Ports))
		for i, p := range proto.Ports {
			if p != nil {
				c.Ports[i] = domain.ContainerPort{
					ContainerPort: int(p.ContainerPort),
					Protocol:      p.Protocol,
				}
			}
		}
	}
	if len(proto.Env) > 0 {
		c.Env = make([]domain.EnvVar, len(proto.Env))
		for i, e := range proto.Env {
			if e != nil {
				c.Env[i] = domain.EnvVar{Name: e.Name, Value: e.Value}
			}
		}
	}
	if proto.Resources != nil {
		c.Resources = &domain.Resources{}
		if proto.Resources.Requests != nil {
			c.Resources.Requests = &domain.ResourceQuantity{
				Memory: proto.Resources.Requests.Memory,
				CPU:    proto.Resources.Requests.Cpu,
			}
		}
		if proto.Resources.Limits != nil {
			c.Resources.Limits = &domain.ResourceQuantity{
				Memory: proto.Resources.Limits.Memory,
				CPU:    proto.Resources.Limits.Cpu,
			}
		}
	}
	if proto.LivenessProbe != nil && proto.LivenessProbe.HttpGet != nil {
		c.LivenessProbe = &domain.Probe{
			HTTPGet:             &domain.HTTPGetAction{Path: proto.LivenessProbe.HttpGet.Path, Port: proto.LivenessProbe.HttpGet.Port},
			InitialDelaySeconds: proto.LivenessProbe.InitialDelaySeconds,
			PeriodSeconds:       proto.LivenessProbe.PeriodSeconds,
		}
	}
	if proto.ReadinessProbe != nil && proto.ReadinessProbe.HttpGet != nil {
		c.ReadinessProbe = &domain.Probe{
			HTTPGet:             &domain.HTTPGetAction{Path: proto.ReadinessProbe.HttpGet.Path, Port: proto.ReadinessProbe.HttpGet.Port},
			InitialDelaySeconds: proto.ReadinessProbe.InitialDelaySeconds,
			PeriodSeconds:       proto.ReadinessProbe.PeriodSeconds,
		}
	}
	return c
}
