package nomad


// NomadClient is an abstraction over the http client  connection 
// to Nomad.
type NomadClient interface {

	HttpClient() (*api.Client, error)
}

// StorageApis provides a means to communicate with Nomad Apis
// w.r.t storage. 
//
// NOTE:
//    A Nomad job spec is treated as a persistent volume storage 
// spec & then submitted to a Nomad deployment. 
//
// NOTE:
//    Nomad has no notion of Persistent Volume.
type StorageApis interface {
  // Create makes a request to Nomad to create a storage resource
  CreateStorage()
  
  // Delete makes a request to Nomad to delete the storage resource
  DeleteStorage()
}

// nomadStorageApi is an implementation of the nomad.StorageApis interface
// This will make API calls to Nomad from mayaserver. In addition, it 
// understands submitting a job specs to a Nomad deployment. 
type nomadStorageApi struct {
}

// Create & submit a job spec that creates a resource in Nomad cluster.
//
// NOTE: 
//    Nomad does not have persistent volume as its first class citizen.
// Hence, this resource should exhibit storage characteristics. The validations
// for this should have been done at the volume plugin implementation.
func (nc *nomadStorageApi) CreateStorage() {
}

// Create & submit a job spec that removes a resource in Nomad cluster.
//
// NOTE: 
//    Nomad does not have persistent volume as its first class citizen.
// Hence, this resource should exhibit storage characteristics. The validations
// for this should have been done at the volume plugin implementation.
func (nc *nomadStorageApi) DeleteStorage() {
}

// Apis provides a means to communicate with Nomad Apis
type Apis interface {

  // This returns a client that can communicate with Nomad
  Client() (NomadClient, error)
  
  // This returns a concrete implementation of StorageApis
  StorageApis() (StorageApis, error)
}

// nomadApiProvider is an implementation of nomad.Apis interface
type nomadApiProvider struct {
}

// Provides a client to communicate with Nomad server.
func (ncp *nomadApiProvider) Client() (NomadClient, error) {
  return &nomadClientUtil{}, nil
}

// Returns an instance of StorageApis.
func (ncp *nomadApiProvider) StorageApis() (StorageApis, error) {
  return &nomadStorageApi{}, nil
}
