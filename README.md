# CDEventer

A Tekton custom task and trigger resources to send CDEvents

## Getting Started

Build and deployment of *CDEventer* is managed via [`ko`]https://github.com/google/ko).
To install CDEventer, clone the repo, point `KO_DOCKER_REPO` to your container
registry:

```shell
export KO_DOCKER_REPO=icr.io/cdeventer
ko apply -f config
```

## Send a CDEvent by creating a `Run`

*CDEventer* looks for `Run` resources in the cluster that reference a resource
of kind `CDEvent` and api version `custom.tekton.dev/v0`. The `CDEvent` resource
does not actually need to exist, the `Run` and its parameters are all that is
required. For example:

```yaml
apiVersion: tekton.dev/v1alpha1
kind: Run
metadata:
  generateName: cdevent-
spec:
  ref:
    apiVersion: custom.tekton.dev/v0
    kind: CDEvent
  params:
    - name: context
      value:
        type: dev.cdevents.taskrun.started.v1
        source: cdeventerRun
    - name: subject
      value:
        id: myTaskRun123
        pipelineName: myPipeline
        url: http://example.com/myTaskRun123
    - name: data
      value:
        customDataContentType: "application/json"
        customData: "{\"k1\": \"v1\"}"
```
