// This file acts a bi-directional plug to below resources:
//
//    1. Mayaserver's orchprovider types
//    2. Hashicorp Nomad's types
package nomad

import (
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/hashicorp/nomad/api"
	"github.com/openebs/mayaserver/lib/api/v1"
	"github.com/openebs/mayaserver/lib/orchprovider"
)

// Name of this orchestration provider.
const NomadOrchProviderName = "nomad"

// This is invoked at startup.
// TODO
//    Put the exact wording than the word `startup` !!
//
// TODO
//    Build a mechanism to reload the orchestrator s.t.
// orchestrator's config can be updated at runtime &
// expected effects can be seen
//
// NOTE:
//    This is a Golang feature.
// Due care needs to be exercised to make sure dependencies are initialized &
// hence available.
func init() {
	orchprovider.RegisterOrchProvider(
		NomadOrchProviderName,
		func(config io.Reader) (orchprovider.OrchestratorInterface, error) {
			return newNomadOrchestrator(config)
		})
}

// NomadOrchestrator is a concrete representation of following
// interfaces:
//
//  1. orchprovider.OrchestratorInterface &
//  2. orchprovider.StoragePlacements
type NomadOrchestrator struct {
	// nStorApis represents an instance capable of invoking
	// storage related APIs
	nStorApis StorageApis

	// nApiClient represents an instance that can make connection &
	// invoke Nomad APIs
	//nApiClient NomadClient

	//region string

	// nConfig represents an instance that provides the coordinates
	// of a Nomad server / cluster deployment.
	nConfig *NomadConfig
}

// newNomadOrchestrator provides a new instance of NomadOrchestrator. This is
// invoked during binary startup.
func newNomadOrchestrator(config io.Reader) (*NomadOrchestrator, error) {

	glog.Infof("Building nomad orchestration provider")

	// Transform the Reader to a NomadConfig
	nCfg, err := readNomadConfig(config)
	if err != nil {
		return nil, fmt.Errorf("unable to read Nomad orchestrator's config: %v", err)
	}

	// TODO
	// validations of the populated config structure

	// Get a new instance of Nomad API Provider
	apis := newNomadApiProvider(nCfg)

	// Get the Nomad api client
	nApiClient, err := apis.Client()
	if err != nil {
		return nil, fmt.Errorf("error creating Nomad api client: %v", err)
	}

	// Get Nomad's storage specific API provider
	nStorApis, err := apis.StorageApis(nApiClient)
	if err != nil {
		return nil, fmt.Errorf("error creating Nomad storage operations instance: %v", err)
	}

	// build the orchestrator instance
	nOrch := &NomadOrchestrator{
		nStorApis: nStorApis,
		//nApiClient: nApiClient,
		nConfig: nCfg,
		//region:   regionName,
	}

	return nOrch, nil
}

// Name provides the name of this orchestrator.
// This is an implementation of the orchprovider.OrchestratorInterface interface.
func (n *NomadOrchestrator) Name() string {

	return NomadOrchProviderName
}

// StoragePlacements is this orchestration provider's
// implementation of the orchprovider.OrchestratorInterface interface.
func (n *NomadOrchestrator) StoragePlacements() (orchprovider.StoragePlacements, bool) {

	return n, true
}

// StorageInfoReq is a contract method implementation of
// orchprovider.StoragePlacements interface. In this implementation,
// a resource details will be fetched from a Nomad deployment.
//
// NOTE:
//    Nomad does not have persistent volume as its first class citizen.
// Hence, this resource should exhibit storage characteristics. The validations
// for this should have been done at the volume plugin implementation.
func (n *NomadOrchestrator) StorageInfoReq(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolume, error) {

	jobName, err := PvcToJobName(pvc)
	if err != nil {
		return nil, err
	}

	job, err := n.nStorApis.StorageInfo(jobName)
	if err != nil {
		return nil, err
	}

	return JobToPv(job)
}

// StoragePlacementReq is a contract method implementation of
// orchprovider.StoragePlacements interface. In this implementation,
// a resource will be created at a Nomad deployment.
//
// NOTE:
//    Nomad does not have persistent volume as its first class citizen.
// Hence, this resource should exhibit storage characteristics. The validations
// for this should have been done at the volume plugin implementation.
func (n *NomadOrchestrator) StoragePlacementReq(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolume, error) {

	// TODO
	// Check for the presence of region
	// If region is present then:
	//    1. fetch the NomadConfig applicable for this region
	//    2. build the StorageApiClient
	//    3. add these to a map in NomadOrchestrator struct
	// If region is not present, then use the available region set against this
	// NomadOrchestrator instance
	// Set the pvc's Label with the region property

	// TODO
	// Check for the presence of datacenter
	// If datacenter is present then fetch its StorageApiClient
	// Set the pvc's Label with the datacenter property

	job, err := PvcToJob(pvc)
	if err != nil {
		return nil, err
	}

	eval, err := n.nStorApis.CreateStorage(job)
	if err != nil {
		return nil, err
	}

	glog.V(2).Infof("Volume '%s' was placed for provisioning with eval '%v'", *job.Name, eval)

	return JobEvalsToPv(*job.Name, []*api.Evaluation{eval})
}

// StorageRemovalReq is a contract method implementation of
// orchprovider.StoragePlacements interface. In this implementation,
// the resource will be removed from the Nomad deployment.
//
// NOTE:
//    Nomad does not have persistent volume as its first class citizen.
// Hence, this resource should exhibit storage characteristics. The validations
// for this should have been done at the volume plugin implementation.
func (n *NomadOrchestrator) StorageRemovalReq(pv *v1.PersistentVolume) (*v1.PersistentVolume, error) {

	job, err := PvToJob(pv)
	if err != nil {
		return nil, err
	}

	eval, err := n.nStorApis.DeleteStorage(job)

	if err != nil {
		return nil, err
	}

	glog.V(2).Infof("Volume '%s' was placed for removal with eval '%v'", pv.Name, eval)

	return JobEvalsToPv(*job.Name, []*api.Evaluation{eval})
}
