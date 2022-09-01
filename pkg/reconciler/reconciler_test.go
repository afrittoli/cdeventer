/*
Copyright 2022 The CDEvents Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
