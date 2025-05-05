package k8sutil_test

import (
	"testing"

	"github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/apis/generic"
	"github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/k8sutil"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestApplyGenericImage(t *testing.T) {
	podTmpl := &corev1.PodTemplateSpec{}
	err := k8sutil.ApplyGenericImage("container", generic.GenericImage{
		Registry:         "docker.io",
		Repo:             "alexandrelamarre",
		Image:            "opentelemetry-stack-gateway",
		Tag:              "v0.1.0",
		ImagePullSecrets: []corev1.LocalObjectReference{},
	}, podTmpl)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(podTmpl.Spec.Containers))
	assert.Equal(t, podTmpl.Spec.Containers[0].Image, "docker.io/alexandrelamarre/opentelemetry-stack-gateway:v0.1.0")
	assert.Equal(t, podTmpl.Spec.ImagePullSecrets, []corev1.LocalObjectReference{})

	podTmpl2 := &corev1.PodTemplateSpec{
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "container",
					Image: "docker.io/alexandrelamarre/opentelemetry-stack-gateway:v0.1.0",
				},
			},
		},
	}

	err = k8sutil.ApplyGenericImage("container", generic.GenericImage{
		Registry: "ghcr.io",
		Repo:     "alexandrelamarre",
		Image:    "gateway",
		Tag:      "v0.1.1",
		ImagePullSecrets: []corev1.LocalObjectReference{
			{
				Name: "test",
			},
		},
	}, podTmpl2)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(podTmpl2.Spec.Containers))
	assert.Equal(t, podTmpl2.Spec.Containers[0].Image, "ghcr.io/alexandrelamarre/gateway:v0.1.1")
	assert.Equal(t, podTmpl2.Spec.ImagePullSecrets, []corev1.LocalObjectReference{
		{
			Name: "test",
		},
	})
}
