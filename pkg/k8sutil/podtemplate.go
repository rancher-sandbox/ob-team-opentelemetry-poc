package k8sutil

import (
	"fmt"
	"strings"

	"github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/apis/generic"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
)

// FIXME : also return a status here
func ApplyGenericImage(container string, image generic.GenericImage, podTmpl *corev1.PodTemplateSpec) error {
	if podTmpl == nil {
		return fmt.Errorf("pod template is nil")
	}

	// modify pod
	podTmpl.Spec.ImagePullSecrets = image.ImagePullSecrets

	// modify containers
	found := false
	for i, cont := range podTmpl.Spec.Containers {
		if cont.Name == container {
			newCont := cont

			imageSb := strings.Builder{}
			if strings.TrimSpace(image.Registry) != "" {
				imageSb.WriteString(image.Registry)
				imageSb.WriteString("/")
			}
			if strings.TrimSpace(image.Repo) != "" {
				imageSb.WriteString(image.Repo)
				imageSb.WriteString("/")
			}
			if strings.TrimSpace(image.Image) != "" {
				imageSb.WriteString(image.Image)
			}
			if strings.TrimSpace(image.Tag) != "" {
				imageSb.WriteString(":")
				imageSb.WriteString(image.Tag)
			}

			// FIXME : dockerparse image as valid, if invalid also return a status

			newCont.Image = imageSb.String()
			podTmpl.Spec.Containers[i] = newCont
			found = true
			break
		}
	}
	if !found {
		imageSb := strings.Builder{}
		if strings.TrimSpace(image.Registry) != "" {
			imageSb.WriteString(image.Registry)
			imageSb.WriteString("/")
		}
		if strings.TrimSpace(image.Repo) != "" {
			imageSb.WriteString(image.Repo)
			imageSb.WriteString("/")
		}
		if strings.TrimSpace(image.Image) != "" {
			imageSb.WriteString(image.Image)
		}
		if strings.TrimSpace(image.Tag) != "" {
			imageSb.WriteString(":")
			imageSb.WriteString(image.Tag)
		}

		podTmpl.Spec.Containers = append(podTmpl.Spec.Containers, corev1.Container{
			Name:  container,
			Image: imageSb.String(),
		})
	}
	return nil
}

func ApplyGenericPodExtension(extension generic.GenericPodExtension, podTmpl *corev1.PodTemplate) error {
	if podTmpl == nil {
		return fmt.Errorf("pod template is nil")
	}
	if extension.ExtraLabels != nil {
		if podTmpl.Labels == nil {
			podTmpl.Labels = map[string]string{}
		}
		podTmpl.Labels = lo.Assign(podTmpl.Labels, extension.ExtraLabels)
	}
	if extension.ExtraAnnotations != nil {
		if podTmpl.Annotations == nil {
			podTmpl.Annotations = map[string]string{}
		}
		podTmpl.Annotations = lo.Assign(podTmpl.Annotations, extension.ExtraAnnotations)
	}
	if extension.ExtraEnv != nil {
		for _, env := range extension.ExtraEnv {
			for i, cont := range podTmpl.Template.Spec.Containers {
				newCont := cont
				skip := false
				for _, contEnv := range cont.Env {
					if contEnv.Name == env.Name {
						skip = true
						break
					}
				}
				if !skip {
					newCont.Env = append(newCont.Env, env)
					podTmpl.Template.Spec.Containers[i] = newCont
				}
			}
		}
	}

	if extension.ExtraVolumes != nil {
		// TODO : dedupe
		podTmpl.Template.Spec.Volumes = append(podTmpl.Template.Spec.Volumes, extension.ExtraVolumes...)
	}

	if extension.ExtraVolumeMounts != nil {
		// TODO
	}

	if extension.ExtraContainers != nil {
		// TODO : dedupe
		podTmpl.Template.Spec.Containers = append(podTmpl.Template.Spec.Containers, extension.ExtraContainers...)
	}

	return nil
}
