package controller

import (
	"github.com/awgreene/status-probe-operator/pkg/controller/probe"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, probe.Add)
}
