package controller

import (
	"github.com/tig4605246/m5g-operator/pkg/controller/mosaic5g"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, mosaic5g.Add)
}
