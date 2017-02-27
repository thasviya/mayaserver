package nomad

// ProviderName is the name of this orchestration provider.
const ProviderName = "nomad"

// NomadConfig provides the settings for NOMAD
type NomadConfig struct {

  Address string
}

// Nomad is a concrete representation of orchprovider.StoragePlacements
// interface
type Nomad struct {
  nStorJobs StorageJobs
  
  nSession  NomadSession
  
  nConfig   *NomadConfig
}

// Services expose all the operations supported by Nomad
type Services interface {

  // This provides a concrete implementation of StorageJobs
  // 
  // NOTE:
  //    Job is a first class citizen in Nomad.
  StorageJobs() (StorageJobs, error)
  
  // This provides a session to communicate with Nomad
  Session() (NomadSession, error)
}

// nomadClientProvider is an implementation of Services interface
type nomadClientProvider struct {

}

// StorageJobs provides the mechanism to re-imagine a Nomad job spec as 
// a persistent volume storage spec & its submission to Nomad
type StorageJobs interface {
  // Create makes a request to Nomad to create a storage resource
  CreateStorageJob()
  
  // Delete makes a request to Nomad to delete the storage resource
  DeleteStorageJob()
}

// nomadClient is an implementation of the StorageJobs interface
type nomadClient struct {
	
}

// NomadSession is an abstraction over the Nomad connection.
type NomadSession interface {

	GetSession() (string, error)
}

// Provides a session to communicate with Nomad server.
func (ncp *nomadClientProvider) Session() (NomadSession, error) {
  return nil, nil
}

// Provides an instance of StorageJobs. This instance of StorageJobs
// understand ways to invoke jobs against a Nomad deployment. 
func (ncp *nomadClientProvider) StorageJobs() (StorageJobs, error) {

  return nil, nil
}

// Create & submit a job spec that creates a resource in Nomad cluster.
// This resource should exhibit storage characteristics.
func (nc *nomadClient) CreateStorageJob() {
}

// Create & submit a job spec that removes a resource in Nomad cluster.
// This removed resource should be storage specific.
func (nc *nomadClient) DeleteStorageJob() {
}

// StoragePlacementReq is this (i.e Nomad) orchestration provider's 
// implementation of the generic orchprovider.StoragePlacements 
// contract. In this implementation, a storage resource will be created at a 
// Nomad deployment.
func (n *Nomad) StoragePlacementReq() {
}

// StorageRemovalReq is this (i.e Nomad) orchestration provider's implementation
// of the generic orchprovider.StoragePlacements contract. In this 
// implementation, a storage resource will be created at a Nomad deployment.
func (n *Nomad) StorageRemovalReq() {
}

