// This file provides the necessary implementation to establish jiva
// as a mayaserver volume plugin.
package jiva

import (
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/openebs/mayaserver/lib/api/v1"
	"github.com/openebs/mayaserver/lib/orchprovider"
	"github.com/openebs/mayaserver/lib/volume"
)

const (
	// A well defined namespaced name given to this volume plugin implementation
	JivaStorPluginName = "openebs.io/jiva"
)

// The registration logic for jiva storage volume plugin
//
// NOTE:
//    This is invoked at startup.
// TODO put the exact word rather than startup !!
//    This is a Golang feature.
//
// NOTE:
//    Due care needs to be exercised to make sure dependencies
// are initialized & hence available.
func init() {
	volume.RegisterVolumePlugin(
		JivaStorPluginName,
		func(config io.Reader, aspect volume.VolumePluginAspect) (volume.VolumeInterface, error) {
			return newJivaStor(config, aspect)
		})
}

// JivaStorNomadAspect is a concrete implementation of VolumePluginAspect
// This is a utility struct (publicly scoped) that can be used during jiva volume
// plugin initialization
type JivaStorNomadAspect struct {

	// The aspect that deals with orchestration needs for jiva
	// storage
	Nomad orchprovider.OrchestratorInterface
}

func (jAspect *JivaStorNomadAspect) GetOrchProvider() (orchprovider.OrchestratorInterface, error) {

	if jAspect.Nomad == nil {
		return nil, fmt.Errorf("Nomad aspect is not set")
	}

	return jAspect.Nomad, nil
}

// jivaStor is the concrete implementation that implements
// following interfaces:
//
//  1. volume.VolumeInterface interface
//  2. volume.Provisioner interface
//  3. volume.Deleter interface
type jivaStor struct {
	// jivaOps abstracts the operations related to this jivaStor
	// instance
	jivaOps JivaOps

	// TODO
	// jConfig provides a handle to tune the operations of
	// this jivaStor instance
	//jConfig *JivaConfig
}

// newJivaStor provides a new instance of jivaStor.
// This function aligns with VolumePluginFactory type.
func newJivaStor(config io.Reader, aspect volume.VolumePluginAspect) (*jivaStor, error) {

	glog.Infof("Building new instance of jiva storage")

	// TODO
	//jCfg, err := readJivaConfig(config)
	//if err != nil {
	//	return nil, fmt.Errorf("unable to read Nomad orchestration provider config file: %v", err)
	//}

	// TODO
	// validations of the populated config structure

	jivaOps, err := newJivaOpsProvider(aspect)
	if err != nil {
		return nil, err
	}

	// build the provisioner instance
	jivaStor := &jivaStor{
		//aspect: aspect,
		jivaOps: jivaOps,
		//jConfig:    jCfg,
	}

	return jivaStor, nil
}

// Name returns the namespaced name of this volume
//
// NOTE:
//    This is a contract implementation of volume.VolumeInterface
func (j *jivaStor) Name() string {
	return JivaStorPluginName
}

// jivaStor supports provisioning
// This is made possible by its jivaOps property
//
// NOTE:
//    This is a contract implementation of volume.VolumeInterface
func (j *jivaStor) Provisioner() (volume.Provisioner, bool) {
	return j, true
}

// jivaStor supports deletion
// This is made possible by its jivaOps property
//
// NOTE:
//    This is a contract implementation of volume.VolumeInterface
func (j *jivaStor) Deleter() (volume.Deleter, bool) {
	return j, true
}

// jivaStor provisions a volume via its jivaOps property.
//
// NOTE:
//    This is a contract implementation of volume.Provisioner interface
func (j *jivaStor) Provision(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolume, error) {

	// TODO
	// Validations of input i.e. claim

	// Delegate to its provider
	pv, err := j.jivaOps.Provision(pvc)

	if err != nil {
		return nil, err
	}

	glog.V(2).Infof("Successfully created jiva volume '%s' '%s'", pv.Name, pv.UID)

	return pv, nil
}

// jivaStor removes a volume via its jivaOps property.
//
// NOTE:
//    This is a contract implementation of volume.Deleter interface
func (j *jivaStor) Delete(pv *v1.PersistentVolume) error {

	// TODO
	// Validations if any

	// Delegate to its provider
	err := j.jivaOps.Delete(pv)

	if err != nil {
		return err
	}

	glog.V(2).Infof("Successfully deleted jiva volume '%s' '%s'", pv.Name, pv.UID)

	return nil

}
