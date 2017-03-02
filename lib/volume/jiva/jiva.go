package jiva

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/openebs/mayaserver/lib/api/v1"
)

// jiva represents the implementation that aligns to volume.Volume
// interface. jiva volumes are disk resources provided by OpenEBS.
//
// NOTE:
//    This will be the base or common struct that can be embedded by
// various action based jiva structures.
type jiva struct {

	// An already created instance of jiva volume
	// TODO
	// Remove this to specific structs if this is not generic
	pv *v1.PersistentVolume

	// Interface that facilitates interaction with jiva provider
	// i.e. orchestrator. The orchestrator related function calls
	// should be initiated from this instance.
	provider jivaProvider

	// TODO
	// A link to its own plugin
	// This is required to get its orchestrator
	// Then why not inject the Aspect only rather than the whole plugin ?
	plugin *jivaVolumePlugin
}

// jivaDeleter represents the implementation that aligns to volume.Deleter
// interface.
type jivaDeleter struct {
	*jiva
}

func (d *jivaDeleter) Delete() error {

	// TODO
	// Validations if any

	// Delegate to its provider
	err := d.provider.DeleteVolume(d)

	if err != nil {
		// Errorf ?
		glog.V(2).Infof("Error deleting JIVA volume '%s' '%s': %v", d.pv.Name, d.pv.UID, err)
		return err
	}

	glog.V(2).Infof("Successfully deleted JIVA volume '%s' '%s'", d.pv.Name, d.pv.UID)

	return nil

}

// jivaProvisioner represents the implementation that aligns to volume.Provisioner
// interface.
type jivaProvisioner struct {
	*jiva

	// volume related options tailored into volume.VolumePluginOptions type
	pvc *v1.PersistentVolumeClaim
}

func (p *jivaProvisioner) Provision() (*v1.PersistentVolume, error) {

	// TODO
	// Validations of input i.e. claim

	// Delegate to its provider
	pv, err := p.provider.CreateVolume(p)

	if err != nil {
		// How to use Errorf ?
		glog.V(2).Infof("Error creating JIVA volume '%s' '%s': %v", pv.Name, pv.UID, err)
		return nil, err
	}

	glog.V(2).Infof("Successfully created JIVA volume '%s' '%s'", pv.Name, pv.UID)

	return pv, nil
}

// jivaProvider interface sets up the blueprint for various jiva volume
// provisioning operations namely creation, deletion, etc.
//
// NOTE:
//    Jiva volume plugin delegates these operations to its provider.
// Hence, the need for this interface.
type jivaProvider interface {

	// CreateVolume will create a jiva volume. It makes use of orchestrator.
	CreateVolume(provisioner *jivaProvisioner) (*v1.PersistentVolume, error)

	// DeleteVolume will delete a jiva volume. It makes use of orchestrator.
	DeleteVolume(deleter *jivaDeleter) error
}

// JivaOrchestrator is the concrete implementation for jivaProvider interface.
type JivaOrchestrator struct{}

// CreateVolume tries to creates a JIVA volume via an orchestrator
func (jOrch *JivaOrchestrator) CreateVolume(p *jivaProvisioner) (*v1.PersistentVolume, error) {
	orchestrator, err := p.plugin.aspect.GetOrchProvider()
	if err != nil {
		return nil, err
	}

	storageOrchestrator, ok := orchestrator.StoragePlacements()

	if !ok {
		return nil, fmt.Errorf("Orchestrator '%s' does not provide storage services", orchestrator.Name())
	}

	return storageOrchestrator.StoragePlacementReq(p.pvc)
}

// DeleteVolume tries to delete the Jiva volume via an orchestrator
func (jOrch *JivaOrchestrator) DeleteVolume(d *jivaDeleter) error {
	orchestrator, err := d.plugin.aspect.GetOrchProvider()
	if err != nil {
		return err
	}

	storageOrchestrator, ok := orchestrator.StoragePlacements()

	if !ok {
		return fmt.Errorf("Orchestrator '%s' does not provide storage services", orchestrator.Name())
	}

	return storageOrchestrator.StorageRemovalReq(d.jiva.pv)
}
