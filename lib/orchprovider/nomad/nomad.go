package nomad

import (
  "github.com/openebs/mayaserver/lib/orchprovider"
)

// Name of this orchestration provider.
const NomadOrchProviderName = "nomad"

// This is invoked during binary startup.
// NOTE: This is a Golang feature.
func init() {
	orchprovider.RegisterOrchProvider(NomadOrchProviderName, func(config io.Reader) (orchprovider.Interface, error) {
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
	nApiClient NomadClient
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
		nStorApis:      nStorApis,
		nApiClient:     nApiClient,
		nConfig:        nCfg,
		//region:   regionName,
	}

	return nOrch, nil
}

// StoragePlacements is this orchestration provider's
// implementation of the orchprovider.Interface interface.
func (n *NomadOrchestrator) StoragePlacements() (StoragePlacements, bool) {

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

