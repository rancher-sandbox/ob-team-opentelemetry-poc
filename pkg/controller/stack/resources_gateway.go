package controller

import (
	"fmt"

	"github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/apis/v1alpha1"
	"github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/k8sutil"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	configServices   = "service"
	configPipelines  = "pipelines"
	configProcessors = "processors"
	configExporters  = "exporters"
	configExtensions = "extensions"
	configReceivers  = "receivers"
)

type ClusterStackGenerator struct {
	systemNamespace string
	clusterstack    *v1alpha1.OpenTelemetryClusterStack
	managedConfig   map[string]any
}

func NewClusterStackGenerator(systemNamespace string, stack *v1alpha1.OpenTelemetryClusterStack) *ClusterStackGenerator {
	return &ClusterStackGenerator{
		systemNamespace: systemNamespace,
		clusterstack:    stack,
		managedConfig:   map[string]any{},
	}
}

type StackGenerator struct {
	stack         *v1alpha1.OpenTelemetryStack
	managedConfig map[string]any
}

func NewStackGenerator(stack *v1alpha1.OpenTelemetryStack) *StackGenerator {
	return &StackGenerator{
		stack:         stack,
		managedConfig: map[string]any{},
	}
}

func (g *StackGenerator) Objects() ([]runtime.Object, error) {
	if g.stack == nil {
		logrus.Info("stack is nil")
		return []runtime.Object{}, nil
	}

	ret := []runtime.Object{}
	if g.stack.Spec.Gateway.Enabled {
		logrus.Info("gateway enabled")
		gwCfg, err := g.gatewayConfigMap()
		if err != nil {
			return nil, err
		}
		logrus.Info("adding config map")
		ret = append(ret, gwCfg)

		gwDeploy, err := g.gatewayDeployment()
		if err != nil {
			return nil, err
		}
		logrus.Info("adding deploy")
		ret = append(ret, gwDeploy)
		logrus.Info("adding service")
		gwSvc := g.gatewayService()
		ret = append(ret, gwSvc)
	}
	return ret, nil
}

func (g *StackGenerator) configMapMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      g.stack.ObjectMeta.Name + "-gateway-config",
		Namespace: g.stack.ObjectMeta.Namespace,
	}
}

func (g *StackGenerator) gatewayConfigMap() (*corev1.ConfigMap, error) {
	if err := g.constructOpenTelemetryConfig(); err != nil {
		return nil, err
	}

	gatewayConfig, err := yaml.Marshal(g.managedConfig)
	if err != nil {
		return nil, err
	}

	cfgmap := &corev1.ConfigMap{
		ObjectMeta: g.configMapMeta(),
		Data: map[string]string{
			"config.yaml": string(gatewayConfig),
		},
	}
	return cfgmap, nil
}

func (g *StackGenerator) constructOpenTelemetryConfig() error {
	if err := g.constructExtensions(); err != nil {
		return fmt.Errorf("failed to construct extensions : %w", err)
	}

	if err := g.constructProcessors(); err != nil {
		return fmt.Errorf("failed to construct processors : %w", err)
	}

	if err := g.constructReceivers(); err != nil {
		return fmt.Errorf("failed to construct receivers : %w", err)
	}

	if err := g.constructExporters(); err != nil {
		return fmt.Errorf("failed to construct exporters : %w", err)
	}

	if err := g.constructPipelines(); err != nil {
		return fmt.Errorf("failed to construct pipelines")
	}
	return nil
}

func (g *StackGenerator) constructProcessors() error {
	proc, err := yamlMarshalString(gatewayProcessors)
	if err != nil {
		return err
	}

	g.managedConfig[configProcessors] = proc
	return nil
}

func (g *StackGenerator) constructExtensions() error {
	ext, err := yamlMarshalString(gatewayExtensions)
	if err != nil {
		return err
	}
	g.managedConfig[configExtensions] = ext
	return nil
}

func (g *StackGenerator) constructReceivers() error {
	recv, err := yamlMarshalString(gatewayReceivers)
	if err != nil {
		return err
	}
	g.managedConfig[configReceivers] = recv
	return nil
}

func (g *StackGenerator) constructExporters() error {
	exp, err := yamlMarshalString(gatewayExporters)
	if err != nil {
		return err
	}
	g.managedConfig[configExporters] = exp
	return nil
}

func (g *StackGenerator) constructPipelines() error {
	extensionMap, ok := g.managedConfig[configExtensions].(map[string]any)
	if !ok {
		return fmt.Errorf("internal error, couldn't extract extensions from managed config")
	}

	receiverMap, ok := g.managedConfig[configReceivers].(map[string]any)
	if !ok {
		return fmt.Errorf("internal error, couldn't extract receivers from managed config")
	}

	exporterMap, ok := g.managedConfig[configExporters].(map[string]any)
	if !ok {
		return fmt.Errorf("internal error, couldn't extract exporters from managed config")
	}

	processorMap, ok := g.managedConfig[configProcessors].(map[string]any)
	if !ok {
		return fmt.Errorf("internal error, couldn't extract processors from managed config")
	}

	registeredExtensions := lo.Keys(extensionMap)
	registeredReceivers := lo.Keys(receiverMap)
	registeredExporters := lo.Keys(exporterMap)
	registeredProcessors := lo.Keys(processorMap)
	g.managedConfig[configServices] = map[string]any{
		configExtensions: registeredExtensions,
		configPipelines: map[string]any{
			"logs": map[string]any{
				"receivers":  registeredReceivers,
				"processors": registeredProcessors,
				"exporters":  registeredExporters,
			},
		},
	}
	return nil
}

var (
	grpcDebugEnvVar = []corev1.EnvVar{
		{
			Name:  "GRPC_GO_LOG_VERBOSITY_LEVEL",
			Value: "99",
		},
		{
			Name:  "GRPC_GO_LOG_SEVERITY_LEVEL",
			Value: "denug",
		},
	}
)

func (g *StackGenerator) gatewayGrpcLogEnvVar() []corev1.EnvVar {
	if g.stack.Spec.Gateway.GrpcDebugLogging {
		return grpcDebugEnvVar
	}
	return []corev1.EnvVar{}
}

func (g *StackGenerator) gatewayDeployment() (*appsv1.Deployment, error) {
	configMeta := g.configMapMeta()
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.stack.Name + "-gateway",
			Namespace: g.stack.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: lo.ToPtr(int32(1)),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"otel.io/stack": g.stack.Name + "-gateway",
				},
			},
			MinReadySeconds:         0,
			RevisionHistoryLimit:    lo.ToPtr(int32(0)),
			ProgressDeadlineSeconds: lo.ToPtr(int32(600)),
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      g.stack.Name + "-gateway",
					Namespace: g.stack.Namespace,
					Labels: map[string]string{
						"otel.io/stack": g.stack.Name + "-gateway",
					},
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									DefaultMode: lo.ToPtr(int32(0644)),
									LocalObjectReference: corev1.LocalObjectReference{
										Name: configMeta.Name,
									},
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name: "gateway",
							Command: []string{
								"collector",
							},
							ImagePullPolicy: corev1.PullAlways,
							Args: []string{
								// "sleep",
								// "3600",
								"--config=filereloader:/var/lib/config.yaml",
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/var/lib",
								},
							},
							Env: append(
								[]corev1.EnvVar{
									{
										Name: "HOST_IP",
										ValueFrom: &corev1.EnvVarSource{
											FieldRef: &corev1.ObjectFieldSelector{
												FieldPath: "status.podIP",
											},
										},
									},
								},
								g.gatewayGrpcLogEnvVar()...,
							),
							Ports: []corev1.ContainerPort{
								{
									Name:          "otlp",
									ContainerPort: 4317,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          "otlphttp",
									ContainerPort: 4318,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          "pprof",
									ContainerPort: 1777,
									Protocol:      corev1.ProtocolTCP,
								},
							},
						},
					},
				},
			},
		},
	}
	if err := k8sutil.ApplyGenericImage("gateway", g.stack.Spec.Gateway.Image, &deploy.Spec.Template); err != nil {
		return nil, err
	}

	return deploy, nil
}

func gatewaySvcMeta(name, namespace string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name + "-gateway",
		Namespace: namespace,
	}
}

func (g *StackGenerator) gatewayService() *corev1.Service {
	svcMeta := gatewaySvcMeta(g.stack.Name, g.stack.Namespace)
	svcMeta.Labels = map[string]string{
		"otel.io/stack": g.stack.Name + "-gateway",
	}
	svc := &corev1.Service{
		ObjectMeta: svcMeta,
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"otel.io/stack": g.stack.Name + "-gateway",
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "otlp",
					Port:       4317,
					TargetPort: intstr.FromString("otlp"),
				},
				{
					Name:       "otlphttp",
					Port:       4318,
					TargetPort: intstr.FromString("otlphttp"),
				},
				{
					Name:       "pprof",
					Port:       1777,
					TargetPort: intstr.FromString("pprof"),
				},
			},
		},
	}

	return svc
}

const serviceFmt = "%s.%s.svc.cluster.local"

const gatewayReceivers = `
otlp:
  protocols:
    grpc:
      endpoint: 0.0.0.0:4317
    http:
      endpoint: 0.0.0.0:4318
`

const gatewayProcessors = `batch:`

const gatewayExtensions = `healthcheckv2:
pprof:
  endpoint: 0.0.0.0:1777
basicauth/opensearch:
  client_auth:
    username: admin
    password: ./fetchOpen==404
`

const gatewayExporters = `debug:
opensearch:
  http:
      tls:
        insecure_skip_verify : true
      endpoint: https://opensearch-cluster-master.default.svc.cluster.local:9200
      auth:
        authenticator: basicauth/opensearch
`

const gatewayConfig = `receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318
processors:
  batch:
extensions:
  healthcheckv2:
  pprof:
    endpoint: 0.0.0.0:1777
exporters:
  debug:
service:
  extensions: [healthcheckv2, pprof]
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [debug]
    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [debug]
    logs:
      receivers: [otlp]
      processors: [batch]
      exporters: [debug]
`
