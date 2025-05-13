package controller

import (
	"path"

	"github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/k8sutil"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (g *ClusterStackGenerator) Objects() ([]runtime.Object, error) {
	if g.clusterstack == nil {
		return []runtime.Object{}, nil
	}

	cfgMap, err := g.configMap()
	if err != nil {
		return nil, err
	}

	daemonSet, err := g.nodeSet(cfgMap.Name)
	if err != nil {
		return nil, err
	}

	return []runtime.Object{cfgMap, daemonSet}, nil
}

func (g *ClusterStackGenerator) configMap() (*corev1.ConfigMap, error) {
	if err := g.constructOpenTelemetryConfig(); err != nil {
		return nil, err
	}

	data, err := yaml.Marshal(g.managedConfig)
	if err != nil {
		return nil, err
	}
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.clusterstack.Name + "-otel-name",
			Namespace: g.systemNamespace,
			Labels:    map[string]string{
				// TODO
			},
		},
		Data: map[string]string{
			"config.yaml": string(data),
		},
	}, nil
}

func (g *ClusterStackGenerator) nodeSet(configMapRef string) (*appsv1.DaemonSet, error) {
	// nodeCfg := g.nodeConfig()
	daemonset := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.clusterstack.Name + "-otel-node",
			Namespace: g.systemNamespace,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"otel.io/stack": g.clusterstack.Name + "otel-node",
				},
			},
			UpdateStrategy: appsv1.DaemonSetUpdateStrategy{
				Type: appsv1.RollingUpdateDaemonSetStrategyType,
			},
			MinReadySeconds:      0,
			RevisionHistoryLimit: lo.ToPtr(int32(0)),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      g.clusterstack.Name + "-otel-node",
					Namespace: g.systemNamespace,
					Labels: map[string]string{
						"otel.io/stack": g.clusterstack.Name + "otel-node",
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
										Name: configMapRef,
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
						// TODO: the following mounts/volume mounts need to be configurable
						{
							Name: "varlogpods",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/log/pods",
								},
							},
						},
						{
							Name: "journal",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/log/journal",
								},
							},
						},
						{
							Name: "auditlogs",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: path.Dir(g.clusterstack.Spec.AuditLogPath),
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
									MountPath: g.stateStoragePath(),
								},
								{
									Name:      "varlogpods",
									MountPath: "/var/log/pods",
									ReadOnly:  true,
								},
								{
									Name:      "journal",
									MountPath: "/var/log/journal",
									ReadOnly:  true,
								},
								{
									Name:      "auditlogs",
									MountPath: path.Dir(g.clusterstack.Spec.AuditLogPath),
									ReadOnly:  true,
								},
							},
							// Env: append(
							// 	[]corev1.EnvVar{},
							// 	g.nodeGrpcLogEnvVar()...,
							// ),
						},
					},
				},
			},
		},
	}

	if err := k8sutil.ApplyGenericImage("node", g.clusterstack.Spec.Image, &daemonset.Spec.Template); err != nil {
		return nil, err
	}

	return daemonset, nil
}
