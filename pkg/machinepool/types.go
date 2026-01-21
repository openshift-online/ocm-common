package machinepool

import (
	"github.com/openshift-online/ocm-common/pkg/cluster"
)

// MachinePoolSpec represents the specification for a machine pool resource
type MachinePoolSpec struct {
	Name     string           `json:"name"`
	Replicas int              `json:"replicas"`
	Platform cluster.Platform `json:"platform"`
}

// MachinePool represents a machine pool resource in the cluster
type MachinePool struct {
	Spec MachinePoolSpec `json:"spec"`
}

// MachinePools represents a collection of machine pools
type MachinePools struct {
	Items []MachinePool `json:"items"`
}
