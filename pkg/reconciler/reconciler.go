package reconciler

import (
	"context"
	"fmt"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"knative.dev/pkg/logging"
	kreconciler "knative.dev/pkg/reconciler"
)

// Reconciler in
type Reconciler struct {
}

// ReconcileKind implements Interface.ReconcileKind.
func (c *Reconciler) ReconcileKind(ctx context.Context, r *v1alpha1.Run) kreconciler.Event {
	logger := logging.FromContext(ctx)
	logger.Infof("Reconciling %s/%s", r.Namespace, r.Name)

	// Ignore completed waits.
	if r.IsDone() {
		logger.Info("Run is finished, done reconciling")
		return nil
	}

	if r.Spec.Ref == nil ||
		r.Spec.Ref.APIVersion != "custom.tekton.dev/v0" || r.Spec.Ref.Kind != "CDEvent" {
		// This is not a Run we should have been notified about; do nothing.
		return nil
	}
	if r.Spec.Ref.Name != "" {
		r.Status.MarkRunFailed("UnexpectedName", "Found unexpected ref name: %s", r.Spec.Ref.Name)
		return fmt.Errorf("unexpected ref name: %s", r.Spec.Ref.Name)
	}

	r.Status.MarkRunSucceeded("Sent", "CDEvent successfully sent")

	// Don't emit events on nop-reconciliations
	return nil
}
