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
	"encoding/json"
	"errors"
	"fmt"

	cdevents "github.com/afrittoli/cdevents-sdk-go/pkg/api"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"knative.dev/pkg/logging"
	kreconciler "knative.dev/pkg/reconciler"
)

const sinkUrl = "http://broker-ingress.knative-eventing.svc.cluster.local/default/events-broker"

type Reconciler struct {
	cloudEventClient cloudevents.Client
}

func missingParam(r *v1alpha1.Run, param string, field string) error {
	r.Status.MarkRunFailed("MissingParam", "Missing \"%s\" in the %s param", param, field)
	return fmt.Errorf("missing \"%s\" in the %s param", param, field)
}

func unexpectedParam(r *v1alpha1.Run, param string, field string) error {
	r.Status.MarkRunFailed("UnexpectedParam", "Unexpected param \"%s\" in %s ", field, param)
	return fmt.Errorf("unexpected param \"%s\" in %s ", field, param)
}

func setSubjectFields(r *v1alpha1.Run, event cdevents.CDEvent, subject map[string]string) {
	// Set the subject fields
	for field, value := range subject {
		switch field {
		case "id":
			event.SetSubjectId(value)
		case "pipelineName":
			switch v := event.(type) {
			case *cdevents.PipelineRunQueuedEvent:
				v.SetSubjectPipelineName(value)
			case *cdevents.PipelineRunStartedEvent:
				v.SetSubjectPipelineName(value)
			case *cdevents.PipelineRunFinishedEvent:
				v.SetSubjectPipelineName(value)
			default:
				unexpectedParam(r, "subject", field)
			}
		case "taskName":
			switch v := event.(type) {
			case *cdevents.TaskRunStartedEvent:
				v.SetSubjectTaskName(value)
			case *cdevents.TaskRunFinishedEvent:
				v.SetSubjectTaskName(value)
			default:
				unexpectedParam(r, "subject", field)
			}
		case "pipelineRun":
			switch v := event.(type) {
			case *cdevents.TaskRunStartedEvent:
				v.SetSubjectPipelineRun(cdevents.Reference{Id: value})
			case *cdevents.TaskRunFinishedEvent:
				v.SetSubjectPipelineRun(cdevents.Reference{Id: value})
			default:
				unexpectedParam(r, "subject", field)
			}
		case "outcome":
			switch v := event.(type) {
			case *cdevents.TaskRunFinishedEvent:
				v.SetSubjectOutcome(cdevents.TaskRunOutcome(value))
			case *cdevents.PipelineRunFinishedEvent:
				v.SetSubjectOutcome(cdevents.PipelineRunOutcome(value))
			default:
				unexpectedParam(r, "subject", field)
			}
		case "url":
			switch v := event.(type) {
			case *cdevents.PipelineRunQueuedEvent:
				v.SetSubjectUrl(value)
			case *cdevents.PipelineRunStartedEvent:
				v.SetSubjectUrl(value)
			case *cdevents.PipelineRunFinishedEvent:
				v.SetSubjectUrl(value)
			case *cdevents.TaskRunStartedEvent:
				v.SetSubjectUrl(value)
			case *cdevents.TaskRunFinishedEvent:
				v.SetSubjectUrl(value)
			case *cdevents.RepositoryCreatedEvent:
				v.SetSubjectUrl(value)
			case *cdevents.RepositoryModifiedEvent:
				v.SetSubjectUrl(value)
			case *cdevents.RepositoryDeletedEvent:
				v.SetSubjectUrl(value)
			case *cdevents.EnvironmentCreatedEvent:
				v.SetSubjectUrl(value)
			case *cdevents.EnvironmentModifiedEvent:
				v.SetSubjectUrl(value)
			default:
				unexpectedParam(r, "subject", field)
			}
		case "errors":
			switch v := event.(type) {
			case *cdevents.TaskRunFinishedEvent:
				v.SetSubjectErrors(value)
			case *cdevents.PipelineRunFinishedEvent:
				v.SetSubjectErrors(value)
			default:
				unexpectedParam(r, "subject", field)
			}
		case "name":
			switch v := event.(type) {
			case *cdevents.RepositoryCreatedEvent:
				v.SetSubjectName(value)
			case *cdevents.RepositoryModifiedEvent:
				v.SetSubjectName(value)
			case *cdevents.RepositoryDeletedEvent:
				v.SetSubjectName(value)
			case *cdevents.EnvironmentCreatedEvent:
				v.SetSubjectName(value)
			case *cdevents.EnvironmentModifiedEvent:
				v.SetSubjectName(value)
			case *cdevents.EnvironmentDeletedEvent:
				v.SetSubjectName(value)
			default:
				unexpectedParam(r, "subject", field)
			}
		case "owner":
			switch v := event.(type) {
			case *cdevents.RepositoryCreatedEvent:
				v.SetSubjectOwner(value)
			case *cdevents.RepositoryModifiedEvent:
				v.SetSubjectOwner(value)
			case *cdevents.RepositoryDeletedEvent:
				v.SetSubjectOwner(value)
			default:
				unexpectedParam(r, "subject", field)
			}
		case "viewUrl":
			switch v := event.(type) {
			case *cdevents.RepositoryCreatedEvent:
				v.SetSubjectViewUrl(value)
			case *cdevents.RepositoryModifiedEvent:
				v.SetSubjectViewUrl(value)
			case *cdevents.RepositoryDeletedEvent:
				v.SetSubjectViewUrl(value)
			default:
				unexpectedParam(r, "subject", field)
			}
		case "artifactId":
			switch v := event.(type) {
			case *cdevents.BuildFinishedEvent:
				v.SetSubjectArtifactId(value)
			default:
				unexpectedParam(r, "subject", field)
			}
		case "environmentId":
			switch v := event.(type) {
			case *cdevents.ServiceDeployedEvent:
				v.SetSubjectEnvironment(cdevents.Reference{Id: value})
			case *cdevents.ServicePublishedEvent:
				v.SetSubjectEnvironment(cdevents.Reference{Id: value})
			case *cdevents.ServiceRemovedEvent:
				v.SetSubjectEnvironment(cdevents.Reference{Id: value})
			case *cdevents.ServiceRolledbackEvent:
				v.SetSubjectEnvironment(cdevents.Reference{Id: value})
			case *cdevents.ServiceUpgradedEvent:
				v.SetSubjectEnvironment(cdevents.Reference{Id: value})
			default:
				unexpectedParam(r, "subject", field)
			}
		default:
			unexpectedParam(r, "subject", field)
		}
	}
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

	var (
		eventContext    map[string]string
		eventSubject    map[string]string
		eventCustomData map[string]string
	)

	// Extract context and subject from params
	for _, param := range r.Spec.Params {
		if param.Value.Type != v1beta1.ParamTypeObject {
			r.Status.MarkRunFailed("UnexpectedParam", "Unexpected param type %v for param %s", param.Value.Type, param.Name)
			return fmt.Errorf("unexpected param type %v for param %s", param.Value.Type, param.Name)
		}
		switch param.Name {
		case "context":
			eventContext = param.Value.ObjectVal
		case "subject":
			eventSubject = param.Value.ObjectVal
		case "data":
			eventCustomData = param.Value.ObjectVal
		default:
			r.Status.MarkRunFailed("UnexpectedParam", "Unexpected param name %s", param.Name)
			return fmt.Errorf("unexpected param name %s", param.Name)
		}
	}

	// Build the CDEvent by type
	eventType, found := eventContext["type"]
	if !found {
		missingParam(r, "context", "type")
	}
	event, err := cdevents.NewCDEvent(cdevents.CDEventType(eventType))
	if err != nil {
		r.Status.MarkRunFailed("InvalidCDEvent", "Could not create CDEvent: %v", err)
		return fmt.Errorf("could not create CDEvent: %v", err)
	}

	// Set the source and subject id
	source, found := eventContext["source"]
	if !found {
		missingParam(r, "context", "source")
	}
	event.SetSource(source)

	// Set all subject fields
	setSubjectFields(r, event, eventSubject)

	// Set the custom data if any
	if eventCustomData != nil {
		customDataBytes, found := eventCustomData["customData"]
		var customData interface{}
		err := json.Unmarshal([]byte(customDataBytes), &customData)
		if err != nil {
			return fmt.Errorf("could not unmarshal custom data: %v", err)
		}
		if !found {
			missingParam(r, "data", "customData")
		}
		customDataContentType, found := eventCustomData["customDataContentType"]
		if !found {
			customDataContentType = "application/json"
		}
		event.SetCustomData(customDataContentType, customData)
	}

	logger.Info(cdevents.AsJsonString(event))

	// As cloud event
	ce, err := cdevents.AsCloudEvent(event)
	if err != nil {
		r.Status.MarkRunFailed("InvalidFormat", "Could not render as CloudEvent %v", err)
		return fmt.Errorf("could not render as CloudEvent %v", err)
	}

	// Now send the event
	if c.cloudEventClient == nil {
		return errors.New("No cloud events client found in the reconciler")
	}
	ctx = cloudevents.ContextWithTarget(context.Background(), sinkUrl)
	ctx = cloudevents.WithEncodingBinary(ctx)
	if result := c.cloudEventClient.Send(ctx, *ce); cloudevents.IsUndelivered(result) {
		r.Status.MarkRunFailed("SendError", "Could not send the CloudEvent: %v", result)
		return fmt.Errorf("could not send the CloudEvent: %v", result)
	}

	r.Status.MarkRunSucceeded("Sent", "CDEvent successfully sent")

	// Don't emit events on nop-reconciliations
	return nil
}
