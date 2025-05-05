package generic

import (
	corev1 "k8s.io/api/core/v1"
)

type GenericImage struct {
	Registry         string                        `json:"registry,omitempty"`
	Repo             string                        `json:"repo"`
	Image            string                        `json:"image"`
	Tag              string                        `json:"tag"`
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
}

func (g *GenericImage) DeepCopyInto(out *GenericImage) {
	*out = *g
}

type GenericPodExtension struct {
	ExtraLabels      map[string]string `json:"extraLabels,omitempty"`
	ExtraAnnotations map[string]string `json:"extraAnnotations,omitempty"`
	ExtraEnv         []corev1.EnvVar   `json:"extraEnv,omitempty"`
	ExtraVolumes     []corev1.Volume   `json:"extraVolumes,omitempty"`
	// ExtraVolumeMounts is a map of container name to volume mounts
	ExtraVolumeMounts map[string][]corev1.VolumeMount `json:"extraVolumeMounts,omitempty"`
	ExtraContainers   []corev1.Container              `json:"extraContainers,omitempty"`
}
