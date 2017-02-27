package nomad

// ProviderName is the name of this orchestration provider.
const ProviderName = "nomad"

// NomadOrchestrator is a concrete representation of following 
// interfaces:
//
//  1. orchprovider.Interface &
//  2. orchprovider.StoragePlacements
type NomadOrchestrator struct {
  nStorApis StorageApis
  
  nClient   NomadClient
  
  nConfig   *NomadConfig
}

// StoragePlacements is this orchestration provider's 
// implementation of the orchprovider.Interface interface.
func (n *NomadOrchestrator) StoragePlacements() (StoragePlacements, bool){

  return n, true
}

// StoragePlacementReq is a contract method implementation of 
// orchprovider.StoragePlacements interface. In this implementation, 
// a resource will be created at a Nomad deployment.
//
// NOTE: 
//    Nomad does not have persistent volume as its first class citizen.
// Hence, this resource should exhibit storage characteristics. The validations
// for this should have been done at the volume plugin implementation.
func (n *NomadOrchestrator) StoragePlacementReq() {
}

// StorageRemovalReq is a contract method implementation of 
// orchprovider.StoragePlacements interface. In this implementation, 
// the resource will be removed from the Nomad deployment.
//
// NOTE: 
//    Nomad does not have persistent volume as its first class citizen.
// Hence, this resource should exhibit storage characteristics. The validations
// for this should have been done at the volume plugin implementation.
func (n *NomadOrchestrator) StorageRemovalReq() {
}
