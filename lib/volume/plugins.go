// This file exposes volume plugin related contracts.
// Hence, any specific volume plugin implementor will
// implement the logic that aligns to the contracts
// exposed here.
// Some of the plugin based interfaces delegate to
// volume based interfaces to do the actual work.
package volume

import (
  "github.com/openebs/mayaserver/lib/api/v1"
)

// Spec is an internal representation of a volume.  
// All API volume types translate to Spec.
type Spec struct {
	Volume           *v1.Volume
	ReadOnly         bool
}

// Name returns the name of Volume.
func (spec *Spec) Name() string {
	switch {
	case spec.Volume != nil:
		return spec.Volume.Name
	default:
		return ""
	}
}

// NewSpecFromVolume creates an Spec from a supplied 
// v1.Volume type
func NewSpecFromVolume(vs *v1.Volume) *Spec {
	return &Spec{
		Volume: vs,
	}
}

// VolumePluginOptions contains option information that will
// be useful to a volume plugin's operations.
type VolumePluginOptions struct {

	// PVC is reference to the claim that lead to provisioning of a new PV.
	// Provisioners *must* create a PV that would be matched by this PVC,
	// i.e. with required capacity, accessMode, labels matching PVC.Selector and
	// so on.
	PVC *v1.PersistentVolumeClaim

	// Tags to attach to the real volume in the storage provider
	Tags *map[string]string

	// Volume provisioning parameters from StorageClass
	Parameters map[string]string
}

// VolumePlugin is an interface to volume based plugins 
// used by mayaserver. This provides the blueprint to instantiate
// and provides other functions that help in managing these plugins.
type VolumePlugin interface {
	// Init initializes the plugin.  This will be called exactly once
	// before any New* calls are made - implementations of plugins may
	// depend on this.
	Init(aspect VolumePluginAspect) error

	// Name returns the plugin's name.  Plugins should use namespaced names
	// such as "org.com/volume".  The "openebs.io" namespace is
	// reserved for plugins which are bundled with openebs.
	GetPluginName() string

	// GetVolumeName returns the name/ID to uniquely identifying the actual
	// backing device, directory, path, etc. referenced by the specified volume
	// spec. If the plugin does not support the given spec, this returns an error.
	GetVolumeName(spec *Spec) (string, error)

	// CanSupport tests whether the plugin supports a given volume
	// specification from the API.  The spec pointer should be considered
	// const.
	CanSupport(spec *Spec) bool

	// ConstructVolumeSpec constructs a volume spec based on the given volume name
	// and mountPath. The spec may have incomplete information due to limited
	// information from input. This function is used by volume manager to reconstruct
	// volume spec by reading the volume directories from disk
	ConstructVolumeSpec(volumeName, mountPath string) (*Spec, error)
}

// ProvisionableVolumePlugin is an extended interface of VolumePlugin and is
// used to create volumes.
type ProvisionableVolumePlugin interface {

	VolumePlugin
	
	// NewProvisioner creates a new volume.Provisioner which knows how to
	// create PersistentVolumes in accordance with the plugin's underlying
	// storage provider
	NewProvisioner(options VolumePluginOptions) (Provisioner, error)
}

// DeletableVolumePlugin is an extended interface of VolumePlugin and is used
// by persistent volumes that want to be deleted from the storage infrastructure
// after their release from a PersistentVolumeClaim.
type DeletableVolumePlugin interface {

	VolumePlugin
	
	// NewDeleter creates a new volume.Deleter which knows how to delete this
	// resource in accordance with the underlying storage provider after the
	// volume's release from a claim
	NewDeleter(spec *Spec) (Deleter, error)
}

// VolumePluginAspect is an interface that provides a blueprint for plugins
// to cater to their needs that stretches beyond volume related operations.
type VolumePluginAspect interface {

	// Get the suitable orchestration provider.
	// A plugin may be linked with its provider e.g.
	// an orchestration provider like K8s, Nomad, Mesos, etc.
	//
	// Note:
	//    OpenEBS believes in running storage software in containers & hence
	// above examples.
	GetOrchProvider() orchprovider.Interface
}

// VolumePluginTracker tracks registered plugins.
type VolumePluginTracker struct {
	mutex   sync.Mutex
	plugins map[string]VolumePlugin
}

// VolumePluginConfig is how volume plugins receive configuration.  An instance
// specific to the plugin will be passed to the plugin's
// ProbeVolumePlugins(config) func.  Reasonable defaults will be provided by
// the binary hosting the plugins while allowing override of those default
// values.  Those config values are then set to an instance of 
// VolumePluginConfig and passed to the plugin.
//
// Values in VolumeConfig are intended to be relevant to several plugins, but
// not necessarily all plugins.  The preference is to leverage strong typing
// in this struct.  All config items must have a descriptive but non-specific
// name (i.e, RecyclerMinimumTimeout is OK but RecyclerMinimumTimeoutForNFS is
// !OK).  An instance of config will be given directly to the plugin, so
// config names specific to plugins are unneeded and wrongly expose plugins in
// this VolumeConfig struct.
//
// OtherAttributes is a map of string values intended for one-off
// configuration of a plugin or config that is only relevant to a single
// plugin.  All values are passed by string and require interpretation by the
// plugin. Passing config as strings is the least desirable option but can be
// used for truly one-off configuration. The binary should still use strong
// typing for this value when binding CLI values before they are passed as
// strings in OtherAttributes.
type VolumePluginConfig struct {

	// OtherAttributes stores config as strings.  These strings are opaque to
	// the system and only understood by the binary hosting the plugin and the
	// plugin itself.
	OtherAttributes map[string]string
}

// InitPlugins initializes each plugin.  All plugins must have unique names.
// This must be called exactly once before any New* methods are called on any
// plugins.
func (pm *VolumePluginTracker) InitPlugins(plugins []VolumePlugin, aspect VolumePluginAspect) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	if pm.plugins == nil {
		pm.plugins = map[string]VolumePlugin{}
	}

	allErrs := []error{}
	for _, plugin := range plugins {
		name := plugin.GetPluginName()
		if errs := validation.IsQualifiedName(name); len(errs) != 0 {
			allErrs = append(allErrs, fmt.Errorf("volume plugin has invalid name: %q: %s", name, strings.Join(errs, ";")))
			continue
		}

		if _, found := pm.plugins[name]; found {
			allErrs = append(allErrs, fmt.Errorf("volume plugin %q was registered more than once", name))
			continue
		}
		err := plugin.Init(aspect)
		if err != nil {
			glog.Errorf("Failed to load volume plugin %s, error: %s", plugin, err.Error())
			allErrs = append(allErrs, err)
			continue
		}
		pm.plugins[name] = plugin
		glog.V(1).Infof("Loaded volume plugin %q", name)
	}
	return utilerrors.NewAggregate(allErrs)
}

// FindPluginBySpec looks for a plugin that can support a given volume
// specification.  If no plugins can support or more than one plugin can
// support it, return error.
func (pm *VolumePluginTracker) FindPluginBySpec(spec *Spec) (VolumePlugin, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	matches := []string{}
	for k, v := range pm.plugins {
		if v.CanSupport(spec) {
			matches = append(matches, k)
		}
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("no volume plugin matched")
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf("multiple volume plugins matched: %s", strings.Join(matches, ","))
	}
	return pm.plugins[matches[0]], nil
}

// FindPluginByName fetches a plugin by name.  If no plugin
// is found, returns error.
func (pm *VolumePluginTracker) FindPluginByName(name string) (VolumePlugin, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Once we can get rid of legacy names we can reduce this to a map lookup.
	matches := []string{}
	for k, v := range pm.plugins {
		if v.GetPluginName() == name {
			matches = append(matches, k)
		}
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("no volume plugin matched")
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf("multiple volume plugins matched: %s", strings.Join(matches, ","))
	}
	return pm.plugins[matches[0]], nil
}

// FindProvisionablePluginByName fetches a provisionable volume plugin by name.  
// If no plugin is found, returns error.
func (pm *VolumePluginTracker) FindProvisionablePluginByName(name string) (ProvisionableVolumePlugin, error) {
	volumePlugin, err := pm.FindPluginByName(name)
	if err != nil {
		return nil, err
	}
	if provisionableVolumePlugin, ok := volumePlugin.(ProvisionableVolumePlugin); ok {
		return provisionableVolumePlugin, nil
	}
	return nil, fmt.Errorf("no provisionable volume plugin matched")
}

// FindDeletablePluginByName fetches a persistent volume plugin by name. If
// no plugin is found, returns error.
func (pm *VolumePluginTracker) FindDeletablePluginByName(name string) (DeletableVolumePlugin, error) {
	volumePlugin, err := pm.FindPluginByName(name)
	if err != nil {
		return nil, err
	}
	if deletableVolumePlugin, ok := volumePlugin.(DeletableVolumePlugin); ok {
		return deletableVolumePlugin, nil
	}
	return nil, fmt.Errorf("no deletable volume plugin matched")
}
