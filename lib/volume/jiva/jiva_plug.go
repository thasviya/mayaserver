// This file provides the necessary implementation to establish jiva
// as a mayaserver volume plugin.
package jiva

import (
	"github.com/openebs/mayaserver/lib/volume"
)

const (
	// A well defined namespaced name given to jiva volume plugin
	jivaVolumePluginName = "openebs.io/jiva"
)

// jivaVolumePlugin is the concrete implementation that aligns to
// volume.VolumePlugin, volume.DeletableVolumePlugin, ProvisionableVolumePlugin
// interfaces. In other words this bridges jiva with mayaserver's volume plugin
// contracts.
type jivaVolumePlugin struct {
	aspect volume.VolumePluginAspect
}

// Init does the generic initialization of jivaVolumePlugin. It sets the aspect
// of jivaVolumePlugin.
func (plugin *jivaVolumePlugin) Init(aspect volume.VolumePluginAspect) error {
	plugin.aspect = aspect
	return nil
}

// GetPluginName returns the namespaced name of this plugin i.e. jivaVolumePlugin
func (plugin *jivaVolumePlugin) GetPluginName() string {
	return jivaVolumePluginName
}

// TODO
//  Check the naming. Is this some type vs. name ?
// GetVolumeName returns the name of the volume specified in the spec.
func (plugin *jivaVolumePlugin) GetVolumeName(spec *volume.Spec) (string, error) {
	volumeSource, _, err := getVolumeSource(spec)
	if err != nil {
		return "", err
	}

	return volumeSource.VolumeID, nil
}

// CanSupport checks whether the supplied spec belongs to here i.e. jiva volume
// plugin
func (plugin *jivaVolumePlugin) CanSupport(spec *volume.Spec) bool {
	return spec.Volume != nil && spec.Volume.Jiva != nil
}

// jivaVolumePlugin provides a concrete implementation of volume.Deleter interface.
// This deleter instance would manage the deletion of a jiva volume.
func (plugin *jivaVolumePlugin) NewDeleter(spec *volume.Spec) (volume.Deleter, error) {
	return plugin.newDeleterInternal(spec, &JivaOrchestrator{})
}

func (plugin *jivaVolumePlugin) newDeleterInternal(spec *volume.Spec, provider jivaProvider) (volume.Deleter, error) {

	return &jivaDeleter{
		jiva: &jiva{
			volName:  spec.Name(),
			volumeID: spec.Volume.Jiva.VolumeID,
			provider: provider,
			plugin:   plugin,
		}}, nil
}

// jivaVolumePlugin provides a concrete implementation of volume.Provisioner
// interface. This provisoner instance would manage the creation of a new jiva
// volume.
func (plugin *jivaVolumePlugin) NewProvisioner(options volume.VolumePluginOptions) (volume.Provisioner, error) {
	return plugin.newProvisionerInternal(options, &JivaOrchestrator{})
}

func (plugin *awsElasticBlockStorePlugin) newProvisionerInternal(options volume.VolumeOptions, provider jivaProvider) (volume.Provisioner, error) {

	return &jivaProvisioner{
		jiva: &jiva{
			provider: provider,
			plugin:   plugin,
		},
		options: options,
	}, nil
}

// This is the primary entrypoint for jiva volume plugin.
// In-fact this is true for all volume plugins.
func ProbeVolumePlugins() []volume.VolumePlugin {
	return []volume.VolumePlugin{&jivaVolumePlugin{nil}}
}

// TODO
//  Check the naming !!! Is this some type vs. source ?
func getVolumeSource(
	spec *volume.Spec) (*v1.JivaVolumeSource, bool, error) {
	if spec.Volume != nil && spec.Volume.Jiva != nil {
		return spec.Volume.Jiva, spec.Volume.Jiva.ReadOnly, nil
	}

	return nil, false, fmt.Errorf("Spec does not reference any JIVA volume type")
}
