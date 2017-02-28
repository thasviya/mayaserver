package nomad

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
  nApiClient NomadClient
}

// Create & submit a job spec that creates a resource in Nomad cluster.
//
// NOTE:
//    Nomad does not have persistent volume as its first class citizen.
// Hence, this resource should exhibit storage characteristics. The validations
// for this should have been done at the volume plugin implementation.
func (nsApi *nomadStorageApi) CreateStorage() error {

  // TODO
  // These are lot of http API calls.
  //
  // Need to think better !!!
  //  1. Do I need so many calls ?
  //  2. Do I need to invoke Deregister on error ?  
  //  3. What is the meaning of ForceEvaluate ? Do I need it?
  //  4. What to do if Summary returns an in-progress state ?
  
  // func (j *Jobs) Info(jobID string, q *QueryOptions) (*Job, *QueryMeta, error)
  job, qMeta, err := nsApi.Http().Jobs.Info(jobID, qOpts)
  
  if err != nil {
    return err
  }
  
  if job != nil {
    // job exists already
    // check the labels, tags, etc
    // goto summary block
    // or
    // return already exists error
  }
  
  // func (j *Jobs) Validate(job *Job, q *WriteOptions) (*JobValidateResponse, *WriteMeta, error)
  jValRes, valMeta, err := nsApi.Http().Jobs.Validate(job, wOpts)
  
  if err != nil {
    return err
  }
  
  //func (j *Jobs) Register(job *Job, q *WriteOptions) (string, *WriteMeta, error)
  evalID, evalMeta, err := nsApi.Http().Jobs.Register(job, wOpts)
  
  if err != nil {
    return err
  }
  
  // func (j *Jobs) Summary(jobID string, q *QueryOptions) (*JobSummary, *QueryMeta, error)
  jSum, sumMeta, err := nsApi.Http().Jobs.Summary(jobID, qOpts)
  
  if err != nil {
    // Check options like:
    //  1. retry on final err,
    //  2. deregister on final err,
    //  3. return on final err
    return err
  }
  
  return nil
}

// Create & submit a job spec that removes a resource in Nomad cluster.
//
// NOTE:
//    Nomad does not have persistent volume as its first class citizen.
// Hence, this resource should exhibit storage characteristics. The validations
// for this should have been done at the volume plugin implementation.
func (nsApi *nomadStorageApi) DeleteStorage() {
}

// Apis provides a means to communicate with Nomad Apis
type Apis interface {

	// This returns a client that can communicate with Nomad
	Client() (NomadClient, error)

	// This returns a concrete implementation of StorageApis
	StorageApis(nApiClient NomadClient) (StorageApis, error)
}

// nomadApiProvider is an implementation of nomad.Apis interface
type nomadApiProvider struct {
}

// newNomadApiProvider provides a new instance of nomadApiProvider
func newNomadApiProvider() *nomadApiProvider {
	return &nomadApiProvider{}
}

// Provides a concrete implementation of Nomad api client that 
// can invoke Nomad APIs
func (ncp *nomadApiProvider) Client() (NomadClient, error) {
	return &nomadClientUtil{}, nil
}

// Returns an instance of StorageApis.
func (ncp *nomadApiProvider) StorageApis(nApiClient NomadClient) (StorageApis, error) {
	return &nomadStorageApi{
	  nApiClient: nApiClient,
	}, nil
}
