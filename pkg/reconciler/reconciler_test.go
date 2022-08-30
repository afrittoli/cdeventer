package reconciler

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"github.com/tektoncd/pipeline/test/diff"
)

func TestReconcile(t *testing.T) {

	testCases := []struct {
		name string
	}{
		{
			name: "first test",
		}, {
			name: "second test",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			r := &Reconciler{}

			if err := r.ReconcileKind(context.Background(), &v1alpha1.Run{}); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if d := cmp.Diff("expected", "observed"); d != "" {
				t.Errorf("diff: %s", diff.PrintWantGot(d))
			}
		})
	}
}
