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
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"k8s.io/client-go/rest"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/logging"
)

func init() {
	injection.Default.RegisterClient(withCloudEventClient)
}

// ceKey is used to associate the CloudEventClient inside the context.Context
type ceKey struct{}

func withCloudEventClient(ctx context.Context, _ *rest.Config) context.Context {
	return withCloudEventClientWithConfig(ctx)
}

func withCloudEventClientWithConfig(ctx context.Context) context.Context {
	logger := logging.FromContext(ctx)

	var useOnceTransport http.RoundTripper = &http.Transport{
		DisableKeepAlives: true,
	}

	p, err := cloudevents.NewHTTP(cloudevents.WithRoundTripper(useOnceTransport))
	if err != nil {
		logger.Panicf("Error creating the cloudevents http protocol: %s", err)
	}

	cloudEventClient, err := cloudevents.NewClient(p, cloudevents.WithUUIDs(), cloudevents.WithTimeNow())
	if err != nil {
		logger.Panicf("Error creating the cloudevents client: %s", err)
	}

	return context.WithValue(ctx, ceKey{}, cloudEventClient)
}

// Get extracts the cloudEventClient client from the context.
func Get(ctx context.Context) cloudevents.Client {
	untyped := ctx.Value(ceKey{})
	if untyped == nil {
		logging.FromContext(ctx).Errorf(
			"Unable to fetch client from context.")
		return nil
	}
	return untyped.(cloudevents.Client)
}

// ToContext adds the cloud events client to the context
func ToContext(ctx context.Context, cec cloudevents.Client) context.Context {
	return context.WithValue(ctx, ceKey{}, cec)
}
