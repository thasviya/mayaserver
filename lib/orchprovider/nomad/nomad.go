package nomad

import (
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/openebs/mayaserver/lib/api/v1"
	"github.com/openebs/mayaserver/lib/orchprovider"
)

// Name of this orchestration provider.
const NomadOrchProviderName = "nomad"

// This is invoked at startup.
// TODO put the exact wording rather than startup !!
//
// NOTE:
//    This is a Golang feature.
// Due care needs to be exercised to make sure dependencies are initialized &
// hence available.
func init() {
	orchprovider.RegisterOrchProvider(
		NomadOrchProviderName,
		func(config io.Reader) (orchprovider.Interface, error) {
			apis := newNomadApiProvider()
			return newNomadOrchestrator(config, apis)
		})
}

// NomadOrchestrator is a concrete representation of following
// interfaces:
//
//  1. orchprovider.Interface &
//  2. orchprovider.StoragePlacements
type NomadOrchestrator struct {
	// nStorApis represents an instance capable of invoking
	// storage related APIs
	nStorApis StorageApis
	// nApiClient represents an instance that can make connection &
	// invoke Nomad APIs
	//nApiClient NomadClient
	// nConfig represents an instance that provides the coordinates
	// of a Nomad server / cluster deployment.
	nConfig *NomadConfig
}

// newNomadOrchestrator provides a new instance of NomadOrchestrator. This is
// invoked during binary startup.
func newNomadOrchestrator(config io.Reader, apis Apis) (*NomadOrchestrator, error) {

	glog.Infof("Building nomad orchestration provider")

	// nomad api client
	nApiClient, err := apis.Client()
	if err != nil {
		return nil, fmt.Errorf("error creating Nomad api client: %v", err)
	}

	nCfg, err := readNomadConfig(config)
	if err != nil {
		return nil, fmt.Errorf("unable to read Nomad orchestration provider config file: %v", err)
	}

	// TODO
	// validations of the populated config structure

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
// This is an implementation of the orchprovider.Interface interface.
func (n *NomadOrchestrator) Name() string {

	return NomadOrchProviderName
}

// StoragePlacements is this orchestration provider's
// implementation of the orchprovider.Interface interface.
func (n *NomadOrchestrator) StoragePlacements() (orchprovider.StoragePlacements, bool) {

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
func (n *NomadOrchestrator) StoragePlacementReq(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolume, error) {

	job, err := PvcToJob(pvc)
	if err != nil {
		return nil, err
	}

	jSum, err := n.nStorApis.CreateStorage(job)
	if err != nil {
		return nil, err
	}

	return JobSummaryToPv(jSum)
}

// StorageRemovalReq is a contract method implementation of
// orchprovider.StoragePlacements interface. In this implementation,
// the resource will be removed from the Nomad deployment.
//
// NOTE:
//    Nomad does not have persistent volume as its first class citizen.
// Hence, this resource should exhibit storage characteristics. The validations
// for this should have been done at the volume plugin implementation.
func (n *NomadOrchestrator) StorageRemovalReq(pv *v1.PersistentVolume) error {

	// TODO
	job, err := PvToJob(pv)
	if err != nil {
		return err
	}

	evalID, err := n.nStorApis.DeleteStorage(job)

	glog.Infof("Volume removal req with eval ID '%s' placed for pv '%s'", evalID, pv.Name)

	if err != nil {
		return err
	}

	return nil
}
