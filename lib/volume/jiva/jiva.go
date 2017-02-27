package jiva

// jiva represents the implementation that aligns to volume.Volume 
// interface. jiva volumes are disk resources provided by OpenEBS.
//
// NOTE:
//    This will be the base or common struct that can be embedded by
// various action based jiva structures.
type jiva struct {

  // Name of the jiva volume, that can be easily remembered by the 
  // operators
	volName string
	
	// Unique id of the volume, used to find the disk resource 
	// in the provider.
	volumeID aws.KubernetesVolumeID
		
	// Interface that facilitates interaction with jiva provider
	provider jivaProvider
	
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
	return d.provider.DeleteVolume(d)
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

// jivaProvider interface sets up the blueprint for various jiva volume
// provisioning operations namely creation, deletion, etc. 
type jivaProvider interface {

  // CreateVolume will create a jiva volume. It makes use of an instance of 
  // volume.Provisioner interface.
	CreateVolume(provisioner *jivaProvisioner) (volumeID volume.MayaVolumeID, volumeSizeGB int, labels map[string]string, err error)
	
	// DeleteVolume will delete a jiva volume. It makes use of an instance of
	// volume.Deleter interface.
	DeleteVolume(deleter *jivaDeleter) error
}

// JivaOrchestrator is the concrete implementation for jivaProvider interface.
// Jiva has a dependency on an orchestrator as its volume
// provider.
type JivaOrchestrator struct{}

// DeleteVolume will delete the jiva resource from appropriate 
// orchestrator
func (jOrch *JivaOrchestrator) DeleteVolume(d *jivaDeleter) error {
	orchestrator, err := d.jiva.plugin.aspect.GetOrchProvider()
	if err != nil {
		return err
	}

	deleted, err := orchestrator.DeleteDisk(d.volumeID)
	if err != nil {
		glog.V(2).Infof("Error deleting JIVA volume %s: %v", d.volumeID, err)
		return err
	}
	if deleted {
		glog.V(2).Infof("Successfully deleted JIVA volume %s", d.volumeID)
	} else {
		glog.V(2).Infof("Successfully deleted JIVA volume %s (actually already deleted)", d.volumeID)
	}
	return nil
}

// CreateVolume creates a JIVA volume.
// Returns: volumeID, volumeSizeGB, labels, error
func (jOrch *JivaOrchestrator) CreateVolume(c *jivaProvisioner) (volume.MayaVolumeID, int, map[string]string, error) {
	orchestrator, err := c.awsElasticBlockStore.plugin.host.GetOrchProvider()
	if err != nil {
		return "", 0, nil, err
	}

	// AWS volumes don't have Name field, store the name in Name tag
	var tags map[string]string
	if c.options.Tags == nil {
		tags = make(map[string]string)
	} else {
		tags = *c.options.Tags
	}
	tags["Name"] = volume.GenerateVolumeName(c.options.ClusterName, c.options.PVName, 255) // AWS tags can have 255 characters

	capacity := c.options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)]
	requestBytes := capacity.Value()
	// Jiva works with gigabytes, convert to GiB with rounding up
	requestGB := int(volume.RoundUpSize(requestBytes, 1024*1024*1024))
	volumeOptions := &aws.VolumeOptions{
		CapacityGB: requestGB,
		Tags:       tags,
		PVCName:    c.options.PVC.Name,
	}
	// Apply Parameters (case-insensitive). We leave validation of
	// the values to the orchestrator.
	for k, v := range c.options.Parameters {
		switch strings.ToLower(k) {
		case "type":
			volumeOptions.VolumeType = v
		case "zone":
			volumeOptions.AvailabilityZone = v
		case "iopspergb":
			volumeOptions.IOPSPerGB, err = strconv.Atoi(v)
			if err != nil {
				return "", 0, nil, fmt.Errorf("invalid iopsPerGB value %q, must be integer between 1 and 30: %v", v, err)
			}
		case "encrypted":
			volumeOptions.Encrypted, err = strconv.ParseBool(v)
			if err != nil {
				return "", 0, nil, fmt.Errorf("invalid encrypted boolean value %q, must be true or false: %v", v, err)
			}
		case "kmskeyid":
			volumeOptions.KmsKeyId = v
		default:
			return "", 0, nil, fmt.Errorf("invalid option %q for volume plugin %s", k, c.plugin.GetPluginName())
		}
	}

	name, err := orchestrator.CreateDisk(volumeOptions)
	if err != nil {
		glog.V(2).Infof("Error creating JIVA volume: %v", err)
		return "", 0, nil, err
	}
	glog.V(2).Infof("Successfully created JIVA volume %s", name)

	labels, err := orchestrator.GetVolumeLabels(name)
	if err != nil {
		// We don't really want to leak the volume here...
		glog.Errorf("error building labels for new JIVA volume %q: %v", name, err)
	}

	return name, int(requestGB), labels, nil
}
