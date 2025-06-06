// Code generated by codegen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"sync"
	"time"

	v1alpha1 "github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/apis/v1alpha1"
	"github.com/rancher/wrangler/v3/pkg/apply"
	"github.com/rancher/wrangler/v3/pkg/condition"
	"github.com/rancher/wrangler/v3/pkg/generic"
	"github.com/rancher/wrangler/v3/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// OpenTelemetryClusterStackController interface for managing OpenTelemetryClusterStack resources.
type OpenTelemetryClusterStackController interface {
	generic.NonNamespacedControllerInterface[*v1alpha1.OpenTelemetryClusterStack, *v1alpha1.OpenTelemetryClusterStackList]
}

// OpenTelemetryClusterStackClient interface for managing OpenTelemetryClusterStack resources in Kubernetes.
type OpenTelemetryClusterStackClient interface {
	generic.NonNamespacedClientInterface[*v1alpha1.OpenTelemetryClusterStack, *v1alpha1.OpenTelemetryClusterStackList]
}

// OpenTelemetryClusterStackCache interface for retrieving OpenTelemetryClusterStack resources in memory.
type OpenTelemetryClusterStackCache interface {
	generic.NonNamespacedCacheInterface[*v1alpha1.OpenTelemetryClusterStack]
}

// OpenTelemetryClusterStackStatusHandler is executed for every added or modified OpenTelemetryClusterStack. Should return the new status to be updated
type OpenTelemetryClusterStackStatusHandler func(obj *v1alpha1.OpenTelemetryClusterStack, status v1alpha1.ClusterStackStatus) (v1alpha1.ClusterStackStatus, error)

// OpenTelemetryClusterStackGeneratingHandler is the top-level handler that is executed for every OpenTelemetryClusterStack event. It extends OpenTelemetryClusterStackStatusHandler by a returning a slice of child objects to be passed to apply.Apply
type OpenTelemetryClusterStackGeneratingHandler func(obj *v1alpha1.OpenTelemetryClusterStack, status v1alpha1.ClusterStackStatus) ([]runtime.Object, v1alpha1.ClusterStackStatus, error)

// RegisterOpenTelemetryClusterStackStatusHandler configures a OpenTelemetryClusterStackController to execute a OpenTelemetryClusterStackStatusHandler for every events observed.
// If a non-empty condition is provided, it will be updated in the status conditions for every handler execution
func RegisterOpenTelemetryClusterStackStatusHandler(ctx context.Context, controller OpenTelemetryClusterStackController, condition condition.Cond, name string, handler OpenTelemetryClusterStackStatusHandler) {
	statusHandler := &openTelemetryClusterStackStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, generic.FromObjectHandlerToHandler(statusHandler.sync))
}

// RegisterOpenTelemetryClusterStackGeneratingHandler configures a OpenTelemetryClusterStackController to execute a OpenTelemetryClusterStackGeneratingHandler for every events observed, passing the returned objects to the provided apply.Apply.
// If a non-empty condition is provided, it will be updated in the status conditions for every handler execution
func RegisterOpenTelemetryClusterStackGeneratingHandler(ctx context.Context, controller OpenTelemetryClusterStackController, apply apply.Apply,
	condition condition.Cond, name string, handler OpenTelemetryClusterStackGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &openTelemetryClusterStackGeneratingHandler{
		OpenTelemetryClusterStackGeneratingHandler: handler,
		apply: apply,
		name:  name,
		gvk:   controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterOpenTelemetryClusterStackStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type openTelemetryClusterStackStatusHandler struct {
	client    OpenTelemetryClusterStackClient
	condition condition.Cond
	handler   OpenTelemetryClusterStackStatusHandler
}

// sync is executed on every resource addition or modification. Executes the configured handlers and sends the updated status to the Kubernetes API
func (a *openTelemetryClusterStackStatusHandler) sync(key string, obj *v1alpha1.OpenTelemetryClusterStack) (*v1alpha1.OpenTelemetryClusterStack, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		if a.condition != "" {
			// Since status has changed, update the lastUpdatedTime
			a.condition.LastUpdated(&newStatus, time.Now().UTC().Format(time.RFC3339))
		}

		var newErr error
		obj.Status = newStatus
		newObj, newErr := a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
		if newErr == nil {
			obj = newObj
		}
	}
	return obj, err
}

type openTelemetryClusterStackGeneratingHandler struct {
	OpenTelemetryClusterStackGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
	seen  sync.Map
}

// Remove handles the observed deletion of a resource, cascade deleting every associated resource previously applied
func (a *openTelemetryClusterStackGeneratingHandler) Remove(key string, obj *v1alpha1.OpenTelemetryClusterStack) (*v1alpha1.OpenTelemetryClusterStack, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1alpha1.OpenTelemetryClusterStack{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	if a.opts.UniqueApplyForResourceVersion {
		a.seen.Delete(key)
	}

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

// Handle executes the configured OpenTelemetryClusterStackGeneratingHandler and pass the resulting objects to apply.Apply, finally returning the new status of the resource
func (a *openTelemetryClusterStackGeneratingHandler) Handle(obj *v1alpha1.OpenTelemetryClusterStack, status v1alpha1.ClusterStackStatus) (v1alpha1.ClusterStackStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.OpenTelemetryClusterStackGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}
	if !a.isNewResourceVersion(obj) {
		return newStatus, nil
	}

	err = generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
	if err != nil {
		return newStatus, err
	}
	a.storeResourceVersion(obj)
	return newStatus, nil
}

// isNewResourceVersion detects if a specific resource version was already successfully processed.
// Only used if UniqueApplyForResourceVersion is set in generic.GeneratingHandlerOptions
func (a *openTelemetryClusterStackGeneratingHandler) isNewResourceVersion(obj *v1alpha1.OpenTelemetryClusterStack) bool {
	if !a.opts.UniqueApplyForResourceVersion {
		return true
	}

	// Apply once per resource version
	key := obj.Namespace + "/" + obj.Name
	previous, ok := a.seen.Load(key)
	return !ok || previous != obj.ResourceVersion
}

// storeResourceVersion keeps track of the latest resource version of an object for which Apply was executed
// Only used if UniqueApplyForResourceVersion is set in generic.GeneratingHandlerOptions
func (a *openTelemetryClusterStackGeneratingHandler) storeResourceVersion(obj *v1alpha1.OpenTelemetryClusterStack) {
	if !a.opts.UniqueApplyForResourceVersion {
		return
	}

	key := obj.Namespace + "/" + obj.Name
	a.seen.Store(key, obj.ResourceVersion)
}
