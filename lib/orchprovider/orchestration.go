package orchprovider

import (
	"github.com/openebs/mayaserver/lib/api/v1"
)

// Interface is an interface to the actual infrastructure.
// It represents a pluggable mechanism for the ochestration
// provider to invoke operations on the infrastructure
// (managed by an orchestrator).
//
// NOTE:
//    The operations currently supported are related to storage
// only.
type Interface interface {

	// Name of the orchestration provider
	Name() string

	// This is a builder for StoragePlacements interface. Will return
	// false if not supported.
	StoragePlacements() (StoragePlacements, bool)
}

// StoragePlacement provides the blueprint for storage related
// placements, scheduling, etc at the orchestrator end.
type StoragePlacements interface {

	// StoragePlacementReq will try to create storage resource(s) at the
	// infrastructure
	StoragePlacementReq(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolume, error)

	// StorageRemovalReq will try to delete the storage resource(s) at
	// the infrastructure
	StorageRemovalReq(pv *v1.PersistentVolume) error
}
