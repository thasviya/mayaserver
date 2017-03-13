// This file models every volume action in form of an interface.
//
// QUERY: How are these related to plugin based interfaces ?
//    The procedure to create concrete instances of these interfaces
// are exposed via concrete instance(s) of volume plugin(s).
package volume

import (
	"github.com/openebs/mayaserver/lib/api/v1"
)

// Volume represents an entity that is created by any
// storage infrastructure provider.
type VolumeInterface interface {
	Name() string

	// This is a builder for Provisioner interface. Will return
	// false if not supported.
	Provisioner() (Provisioner, bool)

	// This is a builder for Deleter interface. Will return
	// false if not supported.
	Deleter() (Deleter, bool)

	// This is a builder for Informer interface. Will return
	// false if not supported.
	Informer() (Informer, bool)
}

// Informer is an interface that can fetch details of the volume from
// the storage infrastructure.
type Informer interface {
	// Info tries to fetch the details of a claim from the underlying storage
	// system. This method returns PersistentVolume representing the
	// already available storage resource.
	Info(*v1.PersistentVolumeClaim) (*v1.PersistentVolume, error)
}

// Provisioner is an interface that can create the volume as a new resource in
// the storage infrastructure.
type Provisioner interface {
	// Provision tries creating (i.e. claim) a resource in the underlying storage
	// system. This method returns PersistentVolume representing the
	// created storage resource.
	Provision(*v1.PersistentVolumeClaim) (*v1.PersistentVolume, error)
}

// Deleter removes the storage resource from the underlying storage infrastructure.
// Any error returned indicates the volume has failed to be reclaimed. A nil
// return indicates success.
type Deleter interface {
	// Delete removes the allocated resource in the storage system.
	Delete(*v1.PersistentVolume) (*v1.PersistentVolume, error)
}
