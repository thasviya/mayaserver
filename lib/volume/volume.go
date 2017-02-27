// This file models every volume action in form of an interface.
// 
// QUERY: How are these related to plugin based interfaces ?
//    The procedure to create concrete instances of these interfaces
// are exposed via concrete instance(s) of volume plugin(s).
package volume

// Volume represents an entity that is created by any
// storage infrastructure provider.
type Volume interface {
  GetName() string
}

// Provisioner is an interface that can create the volume as a new resource in 
// the storage infrastructure.
type Provisioner interface {
	// Provision creates the resource by allocating the underlying volume in a
	// storage system. This method returns PersistentVolume representing the
	// created storage resource.
	Provision() (*v1.PersistentVolume, error)
}

// Deleter removes the storage resource from the underlying storage infrastructure. 
// Any error returned indicates the volume has failed to be reclaimed. A nil 
// return indicates success.
type Deleter interface {

	Volume

	// Delete removes the allocated resource in the storage system.
	Delete() error
}

