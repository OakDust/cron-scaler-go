package k8s

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"scale-handler/internal/domain"
)

const (
	namespace        = "default"
	timezone         = "Europe/Moscow"
	kedaAPIVersion   = "keda.sh/v1alpha1"
	scaledObjectKind = "ScaledObject"
)

var weekdayToCron = map[string]string{
	"monday": "1", "tuesday": "2", "wednesday": "3", "thursday": "4",
	"friday": "5", "saturday": "6", "sunday": "0",
}

type Reconciler struct {
	clientset *kubernetes.Clientset
	dynamic   dynamic.Interface
	logger    *slog.Logger
}

func NewReconciler(kubeconfigPath string, logger *slog.Logger) (*Reconciler, error) {
	// Определяем путь к kubeconfig
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("KUBECONFIG")
	}

	if kubeconfigPath == "" {
		// Проверяем стандартный путь
		home, err := os.UserHomeDir()
		if err == nil {
			defaultPath := filepath.Join(home, ".kube", "config")
			if _, err := os.Stat(defaultPath); err == nil {
				kubeconfigPath = defaultPath
			}
		}
	}

	var config *rest.Config
	var err error

	if kubeconfigPath != "" {
		// Раскрываем путь
		expandedPath := kubeconfigPath
		if strings.HasPrefix(expandedPath, "~/") {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("failed to get user home directory: %w", err)
			}
			expandedPath = filepath.Join(home, expandedPath[2:])
		}

		// Проверяем существование файла
		if _, err := os.Stat(expandedPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("kubeconfig file does not exist: %s", expandedPath)
		}

		config, err = clientcmd.BuildConfigFromFlags("", expandedPath)
		logger.Info("Using kubeconfig", "path", expandedPath)
	} else {
		config, err = rest.InClusterConfig()
		logger.Info("Using in-cluster config")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	return &Reconciler{
		clientset: clientset,
		dynamic:   dyn,
		logger:    logger,
	}, nil
}

func (r *Reconciler) CreateResources(ctx context.Context, schedule *domain.Schedule) error {
	if schedule.Application == nil || len(schedule.Application.Containers) == 0 {
		r.logger.Warn("No application containers, skipping K8s creation", "id", schedule.ID)
		return nil
	}

	name := schedule.ID
	if err := r.createDeployment(ctx, name, schedule.Application); err != nil {
		return err
	}
	return r.createScaledObject(ctx, name, &schedule.Rules)
}

func (r *Reconciler) UpdateResources(ctx context.Context, schedule *domain.Schedule) error {
	if schedule.Application == nil || len(schedule.Application.Containers) == 0 {
		return r.DeleteResources(ctx, schedule.ID)
	}

	name := schedule.ID
	if err := r.updateDeployment(ctx, name, schedule.Application); err != nil {
		return err
	}
	return r.updateScaledObject(ctx, name, &schedule.Rules)
}

func (r *Reconciler) DeleteResources(ctx context.Context, scheduleID string) error {
	client := r.dynamic.Resource(scaledObjectGVR()).Namespace(namespace)
	if err := client.Delete(ctx, scheduleID, metav1.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
		r.logger.Error("Failed to delete ScaledObject", "id", scheduleID, "error", err)
	}

	return r.clientset.AppsV1().Deployments(namespace).Delete(ctx, scheduleID, metav1.DeleteOptions{})
}

func scaledObjectGVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "keda.sh",
		Version:  "v1alpha1",
		Resource: "scaledobjects",
	}
}

func (r *Reconciler) createDeployment(ctx context.Context, name string, app *domain.Application) error {
	deployment := r.buildDeployment(name, app)
	_, err := r.clientset.AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("create deployment: %w", err)
	}
	r.logger.Info("Created Deployment", "name", name, "namespace", namespace)
	return nil
}

func (r *Reconciler) updateDeployment(ctx context.Context, name string, app *domain.Application) error {
	deployment := r.buildDeployment(name, app)
	_, err := r.clientset.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return r.createDeployment(ctx, name, app)
		}
		return fmt.Errorf("update deployment: %w", err)
	}
	r.logger.Info("Updated Deployment", "name", name)
	return nil
}

func (r *Reconciler) buildDeployment(name string, app *domain.Application) *appsv1.Deployment {
	containers := make([]corev1.Container, len(app.Containers))
	for i, c := range app.Containers {
		containers[i] = r.containerToK8s(c)
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(0),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": name},
				},
				Spec: corev1.PodSpec{
					Containers: containers,
				},
			},
		},
	}
}

func (r *Reconciler) containerToK8s(c domain.Container) corev1.Container {
	cont := corev1.Container{
		Name:  c.Name,
		Image: c.Image,
	}
	if len(c.Ports) > 0 {
		cont.Ports = make([]corev1.ContainerPort, len(c.Ports))
		for i, p := range c.Ports {
			proto := corev1.ProtocolTCP
			if strings.ToUpper(p.Protocol) == "UDP" {
				proto = corev1.ProtocolUDP
			}
			cont.Ports[i] = corev1.ContainerPort{
				ContainerPort: int32(p.ContainerPort),
				Protocol:      proto,
			}
		}
	}
	if len(c.Env) > 0 {
		cont.Env = make([]corev1.EnvVar, len(c.Env))
		for i, e := range c.Env {
			cont.Env[i] = corev1.EnvVar{Name: e.Name, Value: e.Value}
		}
	}
	if c.Resources != nil {
		cont.Resources = corev1.ResourceRequirements{}
		if c.Resources.Requests != nil {
			cont.Resources.Requests = corev1.ResourceList{}
			if c.Resources.Requests.Memory != "" {
				cont.Resources.Requests[corev1.ResourceMemory] = resource.MustParse(c.Resources.Requests.Memory)
			}
			if c.Resources.Requests.CPU != "" {
				cont.Resources.Requests[corev1.ResourceCPU] = resource.MustParse(c.Resources.Requests.CPU)
			}
		}
		if c.Resources.Limits != nil {
			cont.Resources.Limits = corev1.ResourceList{}
			if c.Resources.Limits.Memory != "" {
				cont.Resources.Limits[corev1.ResourceMemory] = resource.MustParse(c.Resources.Limits.Memory)
			}
			if c.Resources.Limits.CPU != "" {
				cont.Resources.Limits[corev1.ResourceCPU] = resource.MustParse(c.Resources.Limits.CPU)
			}
		}
	}
	if c.LivenessProbe != nil && c.LivenessProbe.HTTPGet != nil {
		cont.LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: c.LivenessProbe.HTTPGet.Path,
					Port: intstr.FromInt32(c.LivenessProbe.HTTPGet.Port),
				},
			},
			InitialDelaySeconds: c.LivenessProbe.InitialDelaySeconds,
			PeriodSeconds:       c.LivenessProbe.PeriodSeconds,
		}
	}
	if c.ReadinessProbe != nil && c.ReadinessProbe.HTTPGet != nil {
		cont.ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: c.ReadinessProbe.HTTPGet.Path,
					Port: intstr.FromInt32(c.ReadinessProbe.HTTPGet.Port),
				},
			},
			InitialDelaySeconds: c.ReadinessProbe.InitialDelaySeconds,
			PeriodSeconds:       c.ReadinessProbe.PeriodSeconds,
		}
	}
	return cont
}

func (r *Reconciler) createScaledObject(ctx context.Context, name string, rules *domain.ScheduleRules) error {
	obj := r.buildScaledObject(name, rules)
	client := r.dynamic.Resource(scaledObjectGVR()).Namespace(namespace)
	_, err := client.Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("create ScaledObject: %w", err)
	}
	r.logger.Info("Created ScaledObject", "name", name)
	return nil
}

func (r *Reconciler) updateScaledObject(ctx context.Context, name string, rules *domain.ScheduleRules) error {
	obj := r.buildScaledObject(name, rules)
	client := r.dynamic.Resource(scaledObjectGVR()).Namespace(namespace)
	_, err := client.Update(ctx, obj, metav1.UpdateOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return r.createScaledObject(ctx, name, rules)
		}
		return fmt.Errorf("update ScaledObject: %w", err)
	}
	r.logger.Info("Updated ScaledObject", "name", name)
	return nil
}

func (r *Reconciler) buildScaledObject(name string, rules *domain.ScheduleRules) *unstructured.Unstructured {
	triggers := []map[string]interface{}{}

	for day, ranges := range rules.Weekdays {
		dow, ok := weekdayToCron[strings.ToLower(day)]
		if !ok {
			continue
		}
		for _, tr := range ranges {
			start := timeToCron(tr.From, "*", "*", dow) // weekday: day=*, month=*, dow=1-7
			end := timeToCron(tr.To, "*", "*", dow)
			triggers = append(triggers, map[string]interface{}{
				"type": "cron",
				"metadata": map[string]interface{}{
					"timezone":        timezone,
					"start":           start,
					"end":             end,
					"desiredReplicas": strconv.Itoa(int(tr.Replicas)),
				},
			})
		}
	}

	for dateStr, ranges := range rules.Dates {
		parts := strings.Split(dateStr, "-") // YYYY-MM-DD
		if len(parts) != 3 {
			continue
		}
		day, month := parts[2], parts[1] // day=01, month=01
		for _, tr := range ranges {
			start := timeToCron(tr.From, day, month, "*") // date: day=01, month=01, dow=*
			end := timeToCron(tr.To, day, month, "*")
			triggers = append(triggers, map[string]interface{}{
				"type": "cron",
				"metadata": map[string]interface{}{
					"timezone":        timezone,
					"start":           start,
					"end":             end,
					"desiredReplicas": strconv.Itoa(int(tr.Replicas)),
				},
			})
		}
	}

	if len(triggers) == 0 {
		triggers = append(triggers, map[string]interface{}{
			"type": "cron",
			"metadata": map[string]interface{}{
				"timezone":        timezone,
				"start":           "0 0 * * *",
				"end":             "0 1 * * *",
				"desiredReplicas": "0",
			},
		})
	}

	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": kedaAPIVersion,
			"kind":       scaledObjectKind,
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"scaleTargetRef": map[string]interface{}{
					"name": name,
				},
				"minReplicaCount": 0,
				"maxReplicaCount": 100,
				"cooldownPeriod":  300,
				"triggers":        triggers,
			},
		},
	}
}

func int32Ptr(i int32) *int32 {
	return &i
}

func timeToCron(timeStr, day, month, dow string) string {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return "0 0 * * *"
	}
	// Cron: minute hour day month day-of-week
	minute, hour := parts[1], parts[0]
	return fmt.Sprintf("%s %s %s %s %s", minute, hour, day, month, dow)
}
