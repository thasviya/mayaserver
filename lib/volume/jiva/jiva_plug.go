// This file provides the necessary implementation to establish jiva
// as a mayaserver volume plugin.
package jiva

import (
	"io"

	"github.com/golang/glog"
	"github.com/openebs/mayaserver/lib/api/v1"
	"github.com/openebs/mayaserver/lib/volume"
)

const (
	// A well defined namespaced name given to jiva volume plugin
	jivaVolumePluginName = "openebs.io/jiva"
)

// This is invoked at startup.
// TODO put the exact word rather than startup !!)
//
// NOTE:
//    This is a Golang feature.
// Due care needs to be exercised to make sure dependencies are initialized &
// hence available.
func init() {
	volume.RegisterVolumePlugin(
		jivaVolumePluginName,
		func(config io.Reader, aspect volume.VolumePluginAspect) (volume.VolumePlugin, error) {
			return newJivaVolumePlugin(config, aspect)
		})
}

// jivaVolumePlugin is the concrete implementation that aligns to
// volume.VolumePlugin, volume.DeletableVolumePlugin, ProvisionableVolumePlugin
// interfaces. In other words this bridges jiva with mayaserver's volume plugin
// contracts.
type jivaVolumePlugin struct {
	aspect volume.VolumePluginAspect
}

// newJivaVolumePlugin provides a new instance of Jiva VolumePlugin.
// This function aligns with VolumePluginFactory type.
func newJivaVolumePlugin(config io.Reader, aspect volume.VolumePluginAspect) (*jivaVolumePlugin, error) {

	glog.Infof("Building jiva volume plugin")

	// TODO
	//jCfg, err := readJivaConfig(config)
	//if err != nil {
	//	return nil, fmt.Errorf("unable to read Nomad orchestration provider config file: %v", err)
	//}

	// TODO
	// validations of the populated config structure

	// build the provisioner instance
	jivaVolumePlug := &jivaVolumePlugin{
		aspect: aspect,
		//nConfig:    jCfg,
	}

	return jivaVolumePlug, nil
}

// GetPluginName returns the namespaced name of this plugin i.e. jivaVolumePlugin
// This is a contract implementation of volume.VolumePlugin
func (plugin *jivaVolumePlugin) GetPluginName() string {
	return jivaVolumePluginName
}

// jivaVolumePlugin provides a concrete implementation of volume.Deleter interface.
// This deleter instance would manage the deletion of a jiva volume.
func (plugin *jivaVolumePlugin) NewDeleter(pv *v1.PersistentVolume) (volume.Deleter, error) {
	return plugin.newDeleterInternal(pv, &JivaOrchestrator{})
}

func (plugin *jivaVolumePlugin) newDeleterInternal(pv *v1.PersistentVolume, provider jivaProvider) (volume.Deleter, error) {

	return &jivaDeleter{
		jiva: &jiva{
			pv:       pv,
			provider: provider,
			plugin:   plugin,
		}}, nil
}

// jivaVolumePlugin provides a concrete implementation of volume.Provisioner
// interface. This provisoner instance would manage the creation of a new jiva
// volume.
func (plugin *jivaVolumePlugin) NewProvisioner(pvc *v1.PersistentVolumeClaim) (volume.Provisioner, error) {
	return plugin.newProvisionerInternal(pvc, &JivaOrchestrator{})
}

func (plugin *jivaVolumePlugin) newProvisionerInternal(pvc *v1.PersistentVolumeClaim, provider jivaProvider) (volume.Provisioner, error) {

	return &jivaProvisioner{
		jiva: &jiva{
			provider: provider,
			plugin:   plugin,
		},
		pvc: pvc,
	}, nil
}

// This is the primary entrypoint for jiva volume plugin.
// In-fact this is true for all volume plugins.
func ProbeVolumePlugins() []volume.VolumePlugin {
	return []volume.VolumePlugin{&jivaVolumePlugin{nil}}
}
