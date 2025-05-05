package controller

import (
	"fmt"

	"github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/apis/v1alpha1"
	"github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/k8sutil"
	"github.com/samber/lo"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Generator struct {
	stack *v1alpha1.OpenTelemetryStack
}

func (g *Generator) Objects() ([]runtime.Object, error) {
	if g.stack == nil {
		return []runtime.Object{}, nil
	}

	ret := []runtime.Object{}
	if g.stack.Spec.Gateway.Enabled {
		gwCfg := g.gatewayConfigMap()
		ret = append(ret, gwCfg)

		gwDeploy, err := g.gatewayDeployment()
		if err != nil {
			return nil, err
		}
		ret = append(ret, gwDeploy)

		gwSvc := g.gatewayService()
		ret = append(ret, gwSvc)
	}

	if g.stack.Spec.Node.Enabled {
		nodeCfg := g.nodeConfig()
		ret = append(ret, nodeCfg)

		nodeSvc := g.nodeService()
		ret = append(ret, nodeSvc)

		nodeSet, err := g.nodeSet()
		if err != nil {
			return nil, err
		}
		ret = append(ret, nodeSet)

	}
	return ret, nil
}

func (g *Generator) gatewayConfigMap() *corev1.ConfigMap {
	cfgmap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.stack.ObjectMeta.Name + "-gateway-config",
			Namespace: g.stack.ObjectMeta.Namespace,
		},
		Data: map[string]string{
			"config.yaml": gatewayConfig,
		},
	}
	return cfgmap
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

func (g *Generator) gatewayGrpcLogEnvVar() []corev1.EnvVar {
	if g.stack.Spec.Gateway.GrpcDebugLogging {
		return grpcDebugEnvVar
	}
	return []corev1.EnvVar{}
}

func (g *Generator) nodeGrpcLogEnvVar() []corev1.EnvVar {
	if g.stack.Spec.Node.GrpcDebugLogging {
		return grpcDebugEnvVar
	}
	return []corev1.EnvVar{}
}

func (g *Generator) gatewayDeployment() (*appsv1.Deployment, error) {
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
										Name: g.gatewayConfigMap().Name,
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
								"--config=/var/lib/config.yaml",
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

func (g *Generator) gatewayService() *corev1.Service {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.stack.Name + "-gateway",
			Namespace: g.stack.Namespace,
			Labels: map[string]string{
				"otel.io/stack": g.stack.Name + "-gateway",
			},
		},
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

func (g *Generator) nodeConfig() *corev1.ConfigMap {
	gwSvc := g.gatewayService()
	gatewayDns := fmt.Sprintf(serviceFmt, gwSvc.Name, gwSvc.Namespace)
	cfgmap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.stack.Name + "-node-config",
			Namespace: g.stack.Namespace,
		},
		Data: map[string]string{
			"config.yaml": nodeConfigBase + fmt.Sprintf(nodeExp, gatewayDns),
		},
	}
	return cfgmap
}

func (g *Generator) nodeService() *corev1.Service {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.stack.Name + "-node",
			Namespace: g.stack.Namespace,
			Labels: map[string]string{
				"otel.io/stack": g.stack.Name + "-node",
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"otel.io/stack": g.stack.Name + "-node",
			},
			Ports: []corev1.ServicePort{
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

func (g *Generator) nodeSet() (*appsv1.DaemonSet, error) {
	nodeCfg := g.nodeConfig()
	daemonset := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.stack.Name + "-node",
			Namespace: g.stack.Namespace,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"otel.io/stack": g.stack.Name + "-node",
				},
			},
			UpdateStrategy: appsv1.DaemonSetUpdateStrategy{
				Type: appsv1.RollingUpdateDaemonSetStrategyType,
			},
			MinReadySeconds:      0,
			RevisionHistoryLimit: lo.ToPtr(int32(0)),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      g.stack.Name + "-node",
					Namespace: g.stack.Namespace,
					Labels: map[string]string{
						"otel.io/stack": g.stack.Name + "-node",
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
										Name: nodeCfg.Name,
									},
								},
							},
						},
						{
							Name: "filestorage-extension",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/otel/filestorage",
									Type: lo.ToPtr(corev1.HostPathDirectoryOrCreate),
								},
							},
						},
						{
							Name: "varlogpods",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/log/pods",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:            "node",
							ImagePullPolicy: corev1.PullAlways,
							Command: []string{
								"collector",
							},
							Args: []string{
								"--config=/var/lib/config.yaml",
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "pprof",
									ContainerPort: 1777,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/var/lib",
									ReadOnly:  true,
								},
								{
									Name:      "filestorage-extension",
									MountPath: "/var/otel/filestorage",
								},
								{
									Name:      "varlogpods",
									MountPath: "/var/log/pods",
									ReadOnly:  true,
								},
							},
							Env: append(
								[]corev1.EnvVar{},
								g.nodeGrpcLogEnvVar()...,
							),
						},
					},
				},
			},
		},
	}

	if err := k8sutil.ApplyGenericImage("node", g.stack.Spec.Node.Image, &daemonset.Spec.Template); err != nil {
		return nil, err
	}

	return daemonset, nil
}

const serviceFmt = "%s.%s.svc.cluster.local"

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

const nodeConfigBase = `
receivers:
  filelog/k8s:
    include: [ /var/log/pods/*/*/*.log ]
    exclude: []
    storage: file_storage
    include_file_path: true
    include_file_name: false
    operators:
    # Find out which format is used by kubernetes
    - type: router
      id: get-format
      routes:
      - output: parser-docker
        expr: 'body matches "^\\{"'
      - output: parser-crio
        expr: 'body matches "^[^ Z]+ "'
      - output: parser-containerd
        expr: 'body matches "^[^ Z]+Z"'
        # Parse CRI-O format
    - type: regex_parser
      id: parser-crio
      regex: '^(?P<time>[^ Z]+) (?P<stream>stdout|stderr) (?P<logtag>[^ ]*) ?(?P<log>.*)$'
      output: extract_metadata_from_filepath
      timestamp:
        parse_from: attributes.time
        layout_type: gotime
        layout: '2006-01-02T15:04:05.000000000-07:00'
      # Parse CRI-Containerd format
    - type: regex_parser
      id: parser-containerd
      regex: '^(?P<time>[^ ^Z]+Z) (?P<stream>stdout|stderr) (?P<logtag>[^ ]*) ?(?P<log>.*)$'
      output: extract_metadata_from_filepath
      timestamp:
        parse_from: attributes.time
        layout: '%Y-%m-%dT%H:%M:%S.%LZ'
      # Parse Docker format
    - type: json_parser
      id: parser-docker
      output: extract_metadata_from_filepath
      timestamp:
        parse_from: attributes.time
        layout: '%Y-%m-%dT%H:%M:%S.%LZ'
      # Extract metadata from file path
    - type: regex_parser
      id: extract_metadata_from_filepath
      regex: '^.*\/(?P<namespace>[^_]+)_(?P<pod_name>[^_]+)_((?P<confighash>[a-f0-9]{32})|(?P<uid>[0-9a-f]{8}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{4}\b-[0-9a-f]{12}))\/(?P<container_name>[^\._]+)\/(?P<restart_count>\d+)\.log$'
      parse_from: attributes["log.file.path"]
    - type: remove
      field: attributes["log.file.path"]
    # Move out attributes to Attributes
    - type: move
      id: move-namespace
      from: attributes.namespace
      to: resource["k8s.namespace.name"]
    - type: move
      id: move-pod-name
      from: attributes.pod_name
      to: resource["k8s.pod.name"]
    - type: move
      id: move-container-name
      from: attributes.container_name
      to: resource["k8s.container.name"]
    - type: move
      from: attributes.uid
      to: resource["k8s.pod.uid"]
    - type: move
      from: attributes.confighash
      to: resource["k8s.pod.confighash"]
extensions:
  pprof:
    endpoint: 0.0.0.0:1777
  healthcheckv2:
  file_storage:
    directory: /var/otel/filestorage
    timeout: 1s
`
const nodeExp = `exporters:
  otlp/gateway: 
    endpoint : %s:4317
    tls:
      insecure: true
      insecure_skip_verify : true
    compression : none
processors:
  batch:
service:
  telemetry:
    logs:
      level : debug
  extensions: [healthcheckv2, pprof, file_storage]
  pipelines:
    logs:
      receivers: [filelog/k8s]
      processors: [batch]
      exporters: [otlp/gateway]
`
