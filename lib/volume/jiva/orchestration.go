// This file handles jiva storage logic related to mayaserver's orchestration
// provider.
//
// NOTE:
//    jiva storage delegates the provisioning, placement & other operational
// aspects to an orchestration provider. Some of the orchestration providers
// can be Kubernetes, Nomad, etc.
package jiva

import (
	"fmt"

	"github.com/openebs/mayaserver/lib/api/v1"
	"github.com/openebs/mayaserver/lib/volume"
)

type JivaOps interface {
	Info(*v1.PersistentVolumeClaim) (*v1.PersistentVolume, error)

	Provision(*v1.PersistentVolumeClaim) (*v1.PersistentVolume, error)

	Delete(*v1.PersistentVolume) (*v1.PersistentVolume, error)
}

func newJivaOrchestrator(aspect volume.VolumePluginAspect) (JivaOps, error) {
	if aspect == nil {
		return nil, fmt.Errorf("Nil volume plugin aspect was provided")
	}

	return &jivaOrchestrator{
		aspect: aspect,
	}, nil
}

// jivaOrchestrator is the concrete implementation for JivaOps interface.
type jivaOrchestrator struct {
	// Orthogonal concerns and their management w.r.t jiva storage
	// is done via aspect
	aspect volume.VolumePluginAspect
}

// Info tries to fetch details of a jiva volume placed in an orchestrator
func (jOrch *jivaOrchestrator) Info(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolume, error) {

	orchestrator, err := jOrch.aspect.GetOrchProvider()
	if err != nil {
		return nil, err
	}

	storageOrchestrator, ok := orchestrator.StoragePlacements()

	if !ok {
		return nil, fmt.Errorf("Orchestrator '%s' does not provide storage services", orchestrator.Name())
	}

	return storageOrchestrator.StorageInfoReq(pvc)
}

// Provision tries to creates a jiva volume via an orchestrator
func (jOrch *jivaOrchestrator) Provision(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolume, error) {
	orchestrator, err := jOrch.aspect.GetOrchProvider()
	if err != nil {
		return nil, err
	}

	storageOrchestrator, ok := orchestrator.StoragePlacements()

	if !ok {
		return nil, fmt.Errorf("Orchestrator '%s' does not provide storage services", orchestrator.Name())
	}

	return storageOrchestrator.StoragePlacementReq(pvc)
}

// Delete tries to delete the jiva volume via an orchestrator
func (jOrch *jivaOrchestrator) Delete(pv *v1.PersistentVolume) (*v1.PersistentVolume, error) {
	orchestrator, err := jOrch.aspect.GetOrchProvider()
	if err != nil {
		return nil, err
	}

	storageOrchestrator, ok := orchestrator.StoragePlacements()

	if !ok {
		return nil, fmt.Errorf("Orchestrator '%s' does not provide storage services", orchestrator.Name())
	}

	return storageOrchestrator.StorageRemovalReq(pv)
}
