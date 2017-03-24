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

	// NetworkPropsReq will try to fetch the networking details at the orchestrator
	// based on a particular datacenter
	//
	// NetworkPropsReq does not fall under CRUD operations. This is applicable
	// to fetching properties from a config, or database etc.
	//
	// NOTE:
	//    This interface will have no control over Create, Update, Delete operations
	// of network properties
	NetworkPropsReq(dc string) (map[v1.ContainerNetworkingLbl]string, error)
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

	// StoragePropsReq will try to fetch the storage details at the orchestrator
	// based on a particular datacenter
	//
	// StoragePropsReq does not fall under CRUD operations. This is applicable
	// to fetching properties from a config, or database etc.
	//
	// NOTE:
	//    This interface will have no control over Create, Update, Delete operations
	// of storage properties.
	//
	// NOTE:
	//    jiva requires these persistent storage properties to provision
	// its instances e.g. backing persistence location is required on which
	// a jiva replica can operate.
	StoragePropsReq(dc string) (map[v1.ContainerStorageLbl]string, error)
}
