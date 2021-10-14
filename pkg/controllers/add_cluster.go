package controllers

import (
	"github.com/crain-cn/cluster-mesh/pkg/controllers/cluster"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, cluster.Add)
}
