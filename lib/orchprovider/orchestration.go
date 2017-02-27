package orchprovider

// Interface is an interface to the actual infrastructure.
// It represents a pluggable mechanism for the ochestration
// provider to invoke operations on the infrastructure
// (managed by an orchestrator). 
//
// NOTE:
//    The operations currently supported are related to storage 
// only.
type Interface interface {

  // This is a builder for StoragePlacements interface. Will return 
  // false if not supported.
  StoragePlacements() (StoragePlacements, error)
}

// StoragePlacement provides the blueprint for storage related
// placements, scheduling, etc at the orchestrator end.
type StoragePlacements interface {

  // StoragePlacementReq will try to create storage resource(s) at the
  // infrastructure
  StoragePlacementReq()
  
  // StorageRemovalReq will try to delete the storage resource(s) at
  // the infrastructure
  StorageRemovalReq()
}
