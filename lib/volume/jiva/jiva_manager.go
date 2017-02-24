package jiva

// jiva represents the implementation that aligns to volume.Volume 
// interface. jiva volumes are disk resources provided by OpenEBS.
type jiva struct {

  // Name of the jiva volume, that can be easily remembered by the 
  // operators
	volName string
	
	// Unique id of the volume, used to find the disk resource 
	// in the provider.
	volumeID aws.KubernetesVolumeID
		
	// Interface that facilitates interaction with jiva provider
	manager jivaManager
	
	// TODO
	//    Check if this is required ?
	// A link to its own plugin
	plugin  *jivaVolumePlugin	
}

// jivaDeleter represents the implementation that aligns to volume.Deleter
// interface.
type jivaDeleter struct {
	*jiva
}


func (d *jivaDeleter) GetName() string {
	return d.volName
}

func (d *jivaDeleter) Delete() error {
	return d.manager.DeleteVolume(d)
}

// jivaProvisioner represents the implementation that aligns to volume.Provisioner
// interface.
type jivaProvisioner struct {
	*jiva
	
	// volume related options tailored into volume.VolumePluginOptions type
	options   volume.VolumePluginOptions
	
	// TODO
	//    Check if this is required ?
	namespace string
}

// jivaManager interface sets up the blueprint for various jiva volume
// related operations namely creation, deletion, etc. 
//
// NOTE:
//    There is a need for this additional abstraction than just relying on 
// volume.Deleter & volume.Provisioner interfaces is due to the dependency
// of a volume provider. Jiva has a dependency on an orchestrator as its volume
// provider.
type jivaManager interface {

  // CreateVolume will create a jiva volume. It makes use of an instance of 
  // volume.Provisioner interface.
	CreateVolume(provisioner *jivaProvisioner) (volumeID aws.KubernetesVolumeID, volumeSizeGB int, labels map[string]string, err error)
	
	// DeleteVolume will delete a jiva volume. It makes use of an instance of
	// volume.Deleter interface.
	DeleteVolume(deleter *jivaDeleter) error
}

// JivaOrchestrator is the concrete implementation for jivaManager interface.
type JivaOrchestrator struct{}
