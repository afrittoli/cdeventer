package main

import (
	"github.com/afrittoli/cdeventer/pkg/reconciler"
	"knative.dev/pkg/injection/sharedmain"
)

func main() {
	sharedmain.Main(reconciler.ControllerName, reconciler.NewController())
}
