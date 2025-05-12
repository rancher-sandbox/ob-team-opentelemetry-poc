package controller

import (
	"context"

	"github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/apis/v1alpha1"
	"github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/setup"
	"github.com/rancher/wrangler/v3/pkg/apply"
	"github.com/rancher/wrangler/v3/pkg/relatedresource"
	"github.com/samber/lo"
	"k8s.io/apimachinery/pkg/runtime"
)

type handler struct {
	Apply           apply.Apply
	SystemNamespace string
}

func (h *handler) OnClusterStackChange(key string, stack *v1alpha1.OpenTelemetryClusterStack) (*v1alpha1.OpenTelemetryClusterStack, error) {
	applier := h.Apply.WithSetID(
		"node",
	)

	if stack != nil {
		applier = applier.WithOwner(stack)
	}

	g := NewClusterStackGenerator(h.SystemNamespace, stack)

	objs, err := g.Objects()
	if err != nil {
		return stack, nil
	}

	if err := applier.ApplyObjects(objs...); err != nil {
		return stack, err
	}

	return stack, nil
}

func (h *handler) OnStackChange(key string, stack *v1alpha1.OpenTelemetryStack) (*v1alpha1.OpenTelemetryStack, error) {
	applier := h.Apply.WithSetID(
		"gateway",
	)

	if stack != nil {
		applier = applier.WithOwner(stack)
	}

	g := StackGenerator{
		stack: stack,
	}

	objs, err := g.Objects()
	if err != nil {
		return stack, err
	}

	objs = lo.Filter(objs, func(obj runtime.Object, _ int) bool {
		return obj != nil
	})

	if err := applier.ApplyObjects(objs...); err != nil {
		return stack, err
	}
	return stack, nil

}

func (h handler) OnStackRemove(key string, stack *v1alpha1.OpenTelemetryStack) (*v1alpha1.OpenTelemetryStack, error) {
	return stack, nil
}

func Register(ctx context.Context, appCtx *setup.AppContext) {
	// Register the controller
	h := &handler{
		SystemNamespace: "cattle-observability-system",
		Apply: appCtx.Apply.WithCacheTypes(
			appCtx.Core.Service(),
			appCtx.Core.ConfigMap(),
			appCtx.Apps.Deployment(),
			appCtx.Apps.DaemonSet(),
		).WithSetOwnerReference(true, false),
	}
	resolver := relatedresource.OwnerResolver(true, v1alpha1.SchemeGroupVersion.String(), "OpenTelemetryStack")

	relatedresource.Watch(
		appCtx.Context,
		"stack-collector-watch",
		resolver,
		appCtx.Stack,
		appCtx.Core.ConfigMap(),
		appCtx.Core.Service(),
		appCtx.Apps.DaemonSet(),
		appCtx.Apps.Deployment(),
	)
	appCtx.Stack.OnChange(ctx, "stack-controller", h.OnStackChange)
	appCtx.Stack.OnRemove(ctx, "stack-controller", h.OnStackRemove)

}
