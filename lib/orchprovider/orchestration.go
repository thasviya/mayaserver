package orchprovider

import (
	"github.com/openebs/mayaserver/lib/api/v1"
)

// OrchestrationInterface is an interface abstraction of a real orchestrator.
// It represents a pluggable mechanism for any orchestration
// provider to invoke operations on the infrastructure managed by an
// orchestrator.
type OrchestratorInterface interface {

	// Name of the orchestration provider
	Name() string

	// Region where this orchestrator is running/deployed
	Region() string

	// This is a builder for NetworkPlacements interface. Will return
	// false if not supported.
	NetworkPlacements() (NetworkPlacements, bool)

	// This is a builder for StoragePlacements interface. Will return
	// false if not supported.
	StoragePlacements() (StoragePlacements, bool)
}

// NetworkPlacements provides the blueprint for network related
// placements, scheduling, etc at the orchestrator end.
type NetworkPlacements interface {

	// NetworkInfoReq will try to fetch the networking details at the orchestrator
	// based on a particular datacenter
	NetworkInfoReq(dc string) (map[v1.ContainerNetworkingLbl]string, error)
}

// StoragePlacement provides the blueprint for storage related
// placements, scheduling, etc at the orchestrator end.
type StoragePlacements interface {

	// StoragePlacementReq will try to create storage resource(s) at the
	// infrastructure
	StoragePlacementReq(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolume, error)

	// StorageRemovalReq will try to delete the storage resource(s) at
	// the infrastructure
	StorageRemovalReq(pv *v1.PersistentVolume) (*v1.PersistentVolume, error)

	// StorageInfoReq will try to fetch the details of a particular storage
	// resource
	StorageInfoReq(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolume, error)
}
