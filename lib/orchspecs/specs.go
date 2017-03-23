// This is a **Work In Progress** item.
//
// This is being currently developed to create Nomad specs that can be fed
// to Nomad APIs. These specs will be generated based on persistent storage
// properties i.e. scheduling, placement, QoS, etc.
//
// NOTE:
// This will be an effective replacement for lib/orchprovider/nomad/helper_funcs.go.
//
// This will be designed in a way that can cater to K8s specs as well.
package orchspecs

import (
	"fmt"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/helper"
	"github.com/openebs/mayaserver/lib/api/v1"
	v1jiva "github.com/openebs/mayaserver/lib/api/v1/jiva"
)

type Specs interface {
	Name() string

	// This will generate specs/intent that is understood by various
	// orchestrators e.g. Nomad/K8s in order to place storage containers
	NewPlacementSpecs() interface{}
}

// JivaNomadSpecs deals with mapping jiva storage properties against Nomad
// orchestrator's job specs
//
// It implements following interface(s):
//  1.  orchspecs.Specs
type jivaNomadSpecs struct {
	volName string

	volID string

	region string

	dc string

	feVolSize string

	beVolSize string

	feCount int

	beCount int

	beVolumeStor string

	feVersion string

	beVersion string

	networkType string

	feIP string

	beIP string

	feSubnet string

	feInterface string
}

// NewJivaNomadSpecs creates a jivaNomadSpecs instance based on the provided
// persistent volume claim
func NewJivaNomadSpecs(pvc *v1.PersistentVolumeClaim) *jivaNomadSpecs {
	return &jivaNomadSpecs{}
}

func (j *jivaNomadSpecs) Name() string {
	return "JivaNomadSpecs"
}

// NewPlacementSpecs will generate specs/intent that is understood by
// Nomad. This specs determine the scheduling & placement of jiva
// storage based containers
func (j *jivaNomadSpecs) NewPlacementSpecs() interface{} {
	return nil
}

// Transform a jiva based PersistentVolumeClaim type to Nomad job
//
// TODO
// There is redundancy in validation. These should be gone once the transformation
// is handled from jiva namespace.
func JivaPvcToJob(pvc *v1.PersistentVolumeClaim) (*api.Job, error) {

	// series of verifications mandated by jiva based provisioning
	if pvc == nil {
		return nil, fmt.Errorf("Nil jiva persistent volume claim")
	}

	// TODO
	// Aggregate these into one error statement
	if pvc.Name == "" {
		return nil, fmt.Errorf("Name missing in jiva pvc")
	}

	if pvc.Labels == nil {
		return nil, fmt.Errorf("Labels missing in jiva pvc")
	}

	if &pvc.Spec == nil || &pvc.Spec.Resources == nil || pvc.Spec.Resources.Requests == nil {
		return nil, fmt.Errorf("Storage specs missing in jiva pvc")
	}

	if pvc.Labels[string(v1.RegionLbl)] == "" {
		return nil, fmt.Errorf("Region region in jiva pvc")
	}

	if pvc.Labels[string(v1.DatacenterLbl)] == "" {
		return nil, fmt.Errorf("Datacenter missing in jiva pvc")
	}

	if pvc.Labels[string(v1jiva.JivaFrontEndImageLbl)] == "" {
		return nil, fmt.Errorf("FrontEnd image version missing in jiva pvc")
	}

	if pvc.Labels[string(v1.CNTypeLbl)] == "" {
		return nil, fmt.Errorf("CN type missing in jiva pvc")
	}

	if pvc.Labels[string(v1jiva.JivaFrontEndIPLbl)] == "" {
		return nil, fmt.Errorf("FrontEnd ip missing in jiva pvc")
	}

	if pvc.Labels[string(v1jiva.JivaBackEndIPLbl)] == "" {
		return nil, fmt.Errorf("BackEnd ip missing in jiva pvc")
	}

	if pvc.Labels[string(v1.CNSubnetLbl)] == "" {
		return nil, fmt.Errorf("CN subnet missing in jiva pvc")
	}

	if pvc.Labels[string(v1.CNInterfaceLbl)] == "" {
		return nil, fmt.Errorf("CN interface missing in jiva pvc")
	}

	feQuantity := pvc.Spec.Resources.Requests[v1jiva.JivaFrontEndVolSizeLbl]
	feQuantityPtr := &feQuantity

	if feQuantityPtr != nil && feQuantityPtr.Sign() <= 0 {
		return nil, fmt.Errorf("Invalid frontend storage size in jiva pvc")
	}

	beQuantity := pvc.Spec.Resources.Requests[v1jiva.JivaBackEndVolSizeLbl]
	beQuantityPtr := &beQuantity

	if beQuantityPtr != nil && beQuantityPtr.Sign() <= 0 {
		return nil, fmt.Errorf("Invalid backend storage size in jiva pvc")
	}

	jivaFeVolSize := feQuantityPtr.String()
	jivaBeVolSize := beQuantityPtr.String()

	// TODO
	// ID is same as Name currently
	// Do we need to think on it ?
	//jobName := helper.StringToPtr(pvc.Name)
	jobName := pvc.Name
	jobID := pvc.Name
	//region := helper.StringToPtr(pvc.Labels[string(v1.RegionLbl)])
	region := pvc.Labels[string(v1.RegionLbl)]
	dc := pvc.Labels[string(v1.DatacenterLbl)]

	jivaGroupName := "pod"
	jivaVolName := pvc.Name

	// Jiva specific properties that are set during
	// front end task group initialization
	feTaskGroup := "fe" + jivaGroupName
	feTaskGrpCount := 1

	// Jiva specific properties that are set during
	// front end task initialization
	feTaskName := "fe1"
	feTaskDriver := "raw_exec"

	// Jiva specific properties that are set as
	// front end task `Resources`
	feTaskCPU := 500
	feTaskMemMB := 256
	feTaskNetMBits := 400

	// Jiva specific properties that are set as
	// front end task `LogConfig`
	feTaskLogMaxFiles := 3
	feTaskLogFileSizeMB := 1

	// Jiva specific properties that are set as
	// front end task `Config`
	jivaFeCmdLbl := "command"
	jivaFeCmd := "launch-jiva-ctl-with-ip"

	// BE Task Group(s)
	beTaskGroup := "be" + jivaGroupName
	beTaskGrpCount := 1

	// BE Task(s)
	beTaskName := "be1"
	beTaskDriver := "raw_exec"

	// Jiva specific properties that are set as
	// back end task `Resources`
	beTaskCPU := 500
	beTaskMemMB := 256
	beTaskNetMBits := 400

	// Jiva specific properties that are set as
	// back end task `LogConfig`
	beTaskLogMaxFiles := 3
	beTaskLogFileSizeMB := 1

	// Jiva specific properties that are set as
	// back end task `Config`
	jivaBeCmdLbl := "command"
	jivaBeCmd := "launch-jiva-rep-with-ip"

	// Jiva specific properties that are set as task ENV
	jivaFeVersion := pvc.Labels[string(v1jiva.JivaFrontEndImageLbl)]
	jivaBeVersion := pvc.Labels[string(v1jiva.JivaBackEndImageLbl)]
	jivaNetworkType := pvc.Labels[string(v1.CNTypeLbl)]
	jivaFeIP := pvc.Labels[string(v1jiva.JivaFrontEndIPLbl)]
	jivaBeIP := pvc.Labels[string(v1jiva.JivaBackEndIPLbl)]
	jivaFeSubnet := pvc.Labels[string(v1.CNSubnetLbl)]
	jivaFeInterface := pvc.Labels[string(v1.CNInterfaceLbl)]
	jivaBeVolStor := pvc.Labels[string(v1jiva.JivaBackEndVolStor)]

	// Compose front end task group
	feTaskGrp := api.NewTaskGroup(feTaskGroup, feTaskGrpCount)

	// frame the front end restart policy
	feRestartPolicy := &api.RestartPolicy{
		Attempts: helper.IntToPtr(3),
		Interval: helper.TimeToPtr(5 * time.Minute),
		Delay:    helper.TimeToPtr(25 * time.Second),
		Mode:     helper.StringToPtr("delay"),
	}

	// Set restart policy
	feTaskGrp.RestartPolicy = feRestartPolicy

	// Compose the front end task
	feTask := api.NewTask(feTaskName, feTaskDriver).
		SetConfig(jivaFeCmdLbl, jivaFeCmd).
		Require(&api.Resources{
			CPU:      helper.IntToPtr(feTaskCPU),
			MemoryMB: helper.IntToPtr(feTaskMemMB),
			Networks: []*api.NetworkResource{
				&api.NetworkResource{
					MBits: helper.IntToPtr(feTaskNetMBits),
				},
			},
		})

	// Set the log rotation of front end task
	feTask.LogConfig = &api.LogConfig{
		MaxFiles:      helper.IntToPtr(feTaskLogMaxFiles),
		MaxFileSizeMB: helper.IntToPtr(feTaskLogFileSizeMB),
	}

	// Set the front end task ENV properties
	feTask.Env = map[string]string{
		"JIVA_CTL_NAME":    jivaVolName + "-" + feTaskGroup + "-" + feTaskName,
		"JIVA_CTL_VERSION": jivaFeVersion,
		"JIVA_CTL_VOLNAME": jivaVolName,
		"JIVA_CTL_VOLSIZE": jivaFeVolSize,
		"JIVA_CTL_IP":      jivaFeIP,
		"JIVA_CTL_SUBNET":  jivaFeSubnet,
		"JIVA_CTL_IFACE":   jivaFeInterface,
	}

	// Set the front end task Artifacts
	feTask.Artifacts = []*api.TaskArtifact{
		&api.TaskArtifact{
			GetterSource: helper.StringToPtr("https://raw.githubusercontent.com/openebs/jiva/master/scripts/launch-jiva-ctl-with-ip"),
			RelativeDest: helper.StringToPtr("local/"),
		},
	}

	// Add front end task(s) to the task group
	feTaskGrp.AddTask(feTask)

	// Compose back end task group
	beTaskGrp := api.NewTaskGroup(beTaskGroup, beTaskGrpCount)

	// frame the back end restart policy
	beRestartPolicy := &api.RestartPolicy{
		Attempts: helper.IntToPtr(3),
		Interval: helper.TimeToPtr(5 * time.Minute),
		Delay:    helper.TimeToPtr(25 * time.Second),
		Mode:     helper.StringToPtr("delay"),
	}

	// Set restart policy
	beTaskGrp.RestartPolicy = beRestartPolicy

	// Compose the back end task
	beTask := api.NewTask(beTaskName, beTaskDriver).
		SetConfig(jivaBeCmdLbl, jivaBeCmd).
		Require(&api.Resources{
			CPU:      helper.IntToPtr(beTaskCPU),
			MemoryMB: helper.IntToPtr(beTaskMemMB),
			Networks: []*api.NetworkResource{
				&api.NetworkResource{
					MBits: helper.IntToPtr(beTaskNetMBits),
				},
			},
		})

	// Set the log rotation of front end task
	beTask.LogConfig = &api.LogConfig{
		MaxFiles:      helper.IntToPtr(beTaskLogMaxFiles),
		MaxFileSizeMB: helper.IntToPtr(beTaskLogFileSizeMB),
	}

	// Set the front end task ENV properties
	beTask.Env = map[string]string{
		"JIVA_REP_NAME":    jivaVolName + "-" + beTaskGroup + "-" + beTaskName,
		"JIVA_CTL_IP":      jivaFeIP,
		"JIVA_REP_VOLNAME": jivaVolName,
		"JIVA_REP_VOLSIZE": jivaBeVolSize,
		//"JIVA_REP_VOLSTORE": "/tmp/jiva/" + jivaVolName + beTaskGroup + "/" + beTaskName,
		"JIVA_REP_VOLSTORE": jivaBeVolStor + jivaVolName + beTaskGroup + "/" + beTaskName,
		"JIVA_REP_VERSION":  jivaBeVersion,
		"JIVA_REP_NETWORK":  jivaNetworkType,
		"JIVA_REP_IFACE":    jivaFeInterface,
		"JIVA_REP_IP":       jivaBeIP,
		"JIVA_REP_SUBNET":   jivaFeSubnet,
	}

	// Set the front end task Artifacts
	beTask.Artifacts = []*api.TaskArtifact{
		&api.TaskArtifact{
			GetterSource: helper.StringToPtr("https://raw.githubusercontent.com/openebs/jiva/master/scripts/launch-jiva-ctl-with-ip"),
			RelativeDest: helper.StringToPtr("local/"),
		},
	}

	// Add back end task(s) to the task group
	beTaskGrp.AddTask(beTask)

	// Compose a job
	jivaPlacementSpecs := api.NewServiceJob(jobID, jobName, region, 50).
		SetMeta("targetportal", jivaFeIP+":"+v1jiva.JivaIscsiTargetPortalPort).
		SetMeta("iqn", v1jiva.JivaIqnFormatPrefix+":"+jivaVolName).
		AddDatacenter(dc).
		Constrain(api.NewConstraint("${attr.kernel.name}", "=", "linux")).
		AddTaskGroup(feTaskGrp).
		AddTaskGroup(beTaskGrp)

	return jivaPlacementSpecs, nil

}
