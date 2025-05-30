// Code generated by codegen. DO NOT EDIT.

package otel

import (
	v1alpha1 "github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/generated/controllers/otel.stack.io/v1alpha1"
	"github.com/rancher/lasso/pkg/controller"
)

type Interface interface {
	V1alpha1() v1alpha1.Interface
}

type group struct {
	controllerFactory controller.SharedControllerFactory
}

// New returns a new Interface.
func New(controllerFactory controller.SharedControllerFactory) Interface {
	return &group{
		controllerFactory: controllerFactory,
	}
}

func (g *group) V1alpha1() v1alpha1.Interface {
	return v1alpha1.New(g.controllerFactory)
}
