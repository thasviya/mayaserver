// This file exposes volume plugin related contracts.
// Hence, any specific volume plugin implementor will
// implement the logic that aligns to the contracts
// exposed here.
// Some of the plugin based interfaces delegate to
// volume based interfaces to do the actual work.
package volume

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/golang/glog"
	"github.com/openebs/mayaserver/lib/api/v1"
	"github.com/openebs/mayaserver/lib/orchprovider"
)

// VolumePluginFactory is a function that returns a volume.VolumePlugin.
// The config parameter provides an io.Reader handler to the factory in
// order to load specific configurations. If no configuration is provided
// the parameter is nil.
type VolumePluginFactory func(config io.Reader, aspect VolumePluginAspect) (VolumePlugin, error)

// VolumePlugin is an interface to volume based plugins
// used by mayaserver. This provides the blueprint to instantiate
// and provides other functions that help in managing these plugins.
type VolumePlugin interface {
	// Name returns the plugin's name.  Plugins should use namespaced names
	// such as "org.com/volume".  The "openebs.io" namespace is
	// reserved for plugins which are bundled with openebs.
	GetPluginName() string
}

// ProvisionableVolumePlugin is an extended interface of VolumePlugin and is
// used to create volumes.
type ProvisionableVolumePlugin interface {
	VolumePlugin

	// NewProvisioner creates a new volume.Provisioner which knows how to
	// create PersistentVolumes in accordance with the plugin's underlying
	// storage provider
	NewProvisioner(pvc v1.PersistentVolumeClaim) (Provisioner, error)
}

// DeletableVolumePlugin is an extended interface of VolumePlugin and is used
// by persistent volumes that want to be deleted from the storage infrastructure
// after their release from a PersistentVolumeClaim.
type DeletableVolumePlugin interface {
	VolumePlugin

	// NewDeleter creates a new volume.Deleter which knows how to delete this
	// resource in accordance with the underlying storage provider after the
	// volume's release from a claim
	NewDeleter(pv *v1.PersistentVolume) (Deleter, error)
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
	GetOrchProvider() (orchprovider.Interface, error)
}

// All registered volume plugins.
var (
	volumePluginsMutex sync.Mutex

	// A mapped instance of volume plugin name with the plugin's
	// initializer
	volumePlugins = make(map[string]VolumePluginFactory)
)

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

// RegisterVolumePlugin registers a volume.VolumePlugin by name.
//
// NOTE:
//    Each implementation of volume plugin need to call
// RegisterVolumePlugin inside their init() function.
func RegisterVolumePlugin(name string, factory VolumePluginFactory) {
	volumePluginsMutex.Lock()
	defer volumePluginsMutex.Unlock()

	if _, found := volumePlugins[name]; found {
		glog.Fatalf("Volume plugin %q was registered twice", name)
	}

	glog.V(1).Infof("Registered volume plugin %q", name)
	volumePlugins[name] = factory
}

// IsVolumePlugin returns true if name corresponds to an already
// registered volume plugin.
func IsVolumePlugin(name string) bool {
	volumePluginsMutex.Lock()
	defer volumePluginsMutex.Unlock()

	_, found := volumePlugins[name]
	return found
}

// VolumePlugins returns the name of all registered volume
// plugins in a string slice
func VolumePlugins() []string {
	names := []string{}
	volumePluginsMutex.Lock()
	defer volumePluginsMutex.Unlock()

	for name := range volumePlugins {
		names = append(names, name)
	}
	return names
}

// GetVolumePlugin creates an instance of the named volume plugin,
// or nil if the name is unknown. The error return is only used if the named
// volume plugin was known but failed to initialize. The config parameter specifies
// the io.Reader handler of the configuration file for the volume
// plugin, or nil for no configuation.
func GetVolumePlugin(name string, config io.Reader, aspect VolumePluginAspect) (VolumePlugin, error) {
	volumePluginsMutex.Lock()
	defer volumePluginsMutex.Unlock()

	factory, found := volumePlugins[name]
	if !found {
		return nil, nil
	}
	return factory(config, aspect)
}

// TODO
// Who calls this ?
// This will currently be triggered while starting the binary as a http service ?
//
// InitVolumePlugin creates an instance of the named volume plugin.
func InitVolumePlugin(name string, configFilePath string, aspect VolumePluginAspect) (VolumePlugin, error) {
	//var orchestrator Interface
	var volumePlugin VolumePlugin
	var err error

	if name == "" {
		glog.Info("No volume plugin specified.")
		return nil, nil
	}

	if configFilePath != "" {
		var config *os.File
		config, err = os.Open(configFilePath)
		if err != nil {
			glog.Fatalf("Couldn't open volume plugin configuration %s: %#v",
				configFilePath, err)
		}

		defer config.Close()
		volumePlugin, err = GetVolumePlugin(name, config, aspect)
	} else {
		// Pass explicit nil so plugins can actually check for nil. See
		// "Why is my nil error value not equal to nil?" in golang.org/doc/faq.
		volumePlugin, err = GetVolumePlugin(name, nil, aspect)
	}

	if err != nil {
		return nil, fmt.Errorf("could not init volume plugin %q: %v", name, err)
	}

	if volumePlugin == nil {
		return nil, fmt.Errorf("unknown volume plugin %q", name)
	}

	return volumePlugin, nil
}
