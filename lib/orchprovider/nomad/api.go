// This file transforms a Nomad scheduler as an orchestration
// platform for persistent volume placement. OpenEBS calls this as
// placement of storage pod.
package nomad

import (
	"fmt"

	"github.com/hashicorp/nomad/api"
)

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
	CreateStorage(job *api.Job) (*api.JobSummary, error)

	// Delete makes a request to Nomad to delete the storage resource
	DeleteStorage(job *api.Job) (string, error)
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
func (nsApi *nomadStorageApi) CreateStorage(job *api.Job) (*api.JobSummary, error) {

	nApiClient := nsApi.nApiClient
	if nApiClient == nil {
		return nil, fmt.Errorf("nomad api client not initialized")
	}

	nApiHttpClient, err := nApiClient.Http()
	if err != nil {
		return nil, err
	}

	//func (j *Jobs) Register(job *Job, q *WriteOptions) (string, *WriteMeta, error)
	//evalID, evalMeta, err := nApiHttpClient.Jobs().Register(job, &api.WriteOptions{})
	_, _, err = nApiHttpClient.Jobs().Register(job, &api.WriteOptions{})

	if err != nil {
		return nil, err
	}

	// func (j *Jobs) Summary(jobID string, q *QueryOptions) (*JobSummary, *QueryMeta, error)
	//jSum, sumMeta, err := nApiHttpClient.Jobs().Summary(*job.ID, &api.QueryOptions{})
	jSum, _, err := nApiHttpClient.Jobs().Summary(*job.ID, &api.QueryOptions{})

	if err != nil {
		return nil, err
	}

	return jSum, err
}

// Create & submit a job spec that removes a resource in Nomad cluster.
//
// NOTE:
//    Nomad does not have persistent volume as its first class citizen.
// Hence, this resource should exhibit storage characteristics. The validations
// for this should have been done at the volume plugin implementation.
func (nsApi *nomadStorageApi) DeleteStorage(job *api.Job) (string, error) {

	nApiClient := nsApi.nApiClient
	if nApiClient == nil {
		return "", fmt.Errorf("nomad api client not initialized")
	}

	nApiHttpClient, err := nApiClient.Http()
	if err != nil {
		return "", err
	}

	//func (j *Jobs) Deregister(jobID string, q *WriteOptions) (string, *WriteMeta, error)
	//evalID, deregMeta, err := nApiHttpClient.Jobs().Deregister(*job.ID, &api.WriteOptions{})
	evalID, _, err := nApiHttpClient.Jobs().Deregister(*job.ID, &api.WriteOptions{})

	if err != nil {
		return "", err
	}

	return evalID, nil
}
