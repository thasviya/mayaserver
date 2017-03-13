package nomad

import (
	"fmt"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/helper"
	"github.com/openebs/mayaserver/lib/api/v1"
)

// Transform a PersistentVolumeClaim type to Nomad job type
func PvcToJob(pvc *v1.PersistentVolumeClaim) (*api.Job, error) {

	if pvc == nil {
		return nil, fmt.Errorf("Nil persistent volume claim provided")
	}

	if pvc.Name == "" {
		return nil, fmt.Errorf("Missing name in persistent volume claim")
	}

	if pvc.Labels == nil {
		return nil, fmt.Errorf("Missing labels in persistent volume claim")
	}

	if pvc.Labels["region"] == "" {
		return nil, fmt.Errorf("Missing region in persistent volume claim")
	}

	if pvc.Labels["datacenter"] == "" {
		return nil, fmt.Errorf("Missing datacenter in persistent volume claim")
	}

	if pvc.Labels["jivafeversion"] == "" {
		return nil, fmt.Errorf("Missing jiva fe version in persistent volume claim")
	}

	if pvc.Labels["jivafenetwork"] == "" {
		return nil, fmt.Errorf("Missing jiva fe network in persistent volume claim")
	}

	if pvc.Labels["jivafeip"] == "" {
		return nil, fmt.Errorf("Missing jiva fe ip in persistent volume claim")
	}

	if pvc.Labels["jivabeip"] == "" {
		return nil, fmt.Errorf("Missing jiva be ip in persistent volume claim")
	}

	if pvc.Labels["jivafesubnet"] == "" {
		return nil, fmt.Errorf("Missing jiva fe subnet in persistent volume claim")
	}

	if pvc.Labels["jivafeinterface"] == "" {
		return nil, fmt.Errorf("Missing jiva fe interface in persistent volume claim")
	}

	// TODO
	// ID is same as Name currently
	// Do we need to think on it ?
	jobName := helper.StringToPtr(pvc.Name)
	region := helper.StringToPtr(pvc.Labels["region"])
	dc := pvc.Labels["datacenter"]

	jivaGroupName := "pod"
	jivaVolName := pvc.Name
	jivaVolSize := "5g"

	feTaskGroup := "fe" + jivaGroupName
	feTaskName := "fe1"
	beTaskGroup := "be" + jivaGroupName
	beTaskName := "be1"

	jivaFeVersion := pvc.Labels["jivafeversion"]
	jivaFeNetwork := pvc.Labels["jivafenetwork"]
	jivaFeIP := pvc.Labels["jivafeip"]
	jivaBeIP := pvc.Labels["jivabeip"]
	jivaFeSubnet := pvc.Labels["jivafesubnet"]
	jivaFeInterface := pvc.Labels["jivafeinterface"]

	// TODO
	// Transformation from pvc or pv to nomad types & vice-versa:
	//
	//  1. Need an Interface or functional callback defined at
	// lib/api/v1/nomad.go &
	//  2. implemented by the volume plugins that want
	// to be orchestrated by Nomad
	//  3. This transformer instance needs to be injected from
	// volume plugin to orchestrator, in a generic way.

	// Hardcoded logic all the way
	// Nomad specific defaults, hardcoding is OK.
	// However, volume plugin specific stuff is BAD
	return &api.Job{
		Region:      region,
		Name:        jobName,
		ID:          jobName,
		Datacenters: []string{dc},
		Type:        helper.StringToPtr(api.JobTypeService),
		Priority:    helper.IntToPtr(50),
		Constraints: []*api.Constraint{
			api.NewConstraint("${attr.kernel.name}", "=", "linux"),
		},
		TaskGroups: []*api.TaskGroup{
			// jiva frontend
			&api.TaskGroup{
				Name:  helper.StringToPtr(feTaskGroup),
				Count: helper.IntToPtr(1),
				RestartPolicy: &api.RestartPolicy{
					Attempts: helper.IntToPtr(3),
					Interval: helper.TimeToPtr(5 * time.Minute),
					Delay:    helper.TimeToPtr(25 * time.Second),
					Mode:     helper.StringToPtr("delay"),
				},
				Tasks: []*api.Task{
					&api.Task{
						Name:   feTaskName,
						Driver: "raw_exec",
						Resources: &api.Resources{
							CPU:      helper.IntToPtr(500),
							MemoryMB: helper.IntToPtr(256),
							Networks: []*api.NetworkResource{
								&api.NetworkResource{
									MBits: helper.IntToPtr(400),
								},
							},
						},
						Env: map[string]string{
							"JIVA_CTL_NAME":    pvc.Name + "-" + feTaskGroup + "-" + feTaskName,
							"JIVA_CTL_VERSION": jivaFeVersion,
							"JIVA_CTL_VOLNAME": jivaVolName,
							"JIVA_CTL_VOLSIZE": jivaVolSize,
							"JIVA_CTL_IP":      jivaFeIP,
							"JIVA_CTL_SUBNET":  jivaFeSubnet,
							"JIVA_CTL_IFACE":   jivaFeInterface,
						},
						Artifacts: []*api.TaskArtifact{
							&api.TaskArtifact{
								GetterSource: helper.StringToPtr("https://raw.githubusercontent.com/openebs/jiva/master/scripts/launch-jiva-ctl-with-ip"),
								RelativeDest: helper.StringToPtr("local/"),
							},
						},
						Config: map[string]interface{}{
							"command": "launch-jiva-ctl-with-ip",
						},
						LogConfig: &api.LogConfig{
							MaxFiles:      helper.IntToPtr(3),
							MaxFileSizeMB: helper.IntToPtr(1),
						},
					},
				},
			},
			// jiva replica
			&api.TaskGroup{
				Name:  helper.StringToPtr(beTaskGroup),
				Count: helper.IntToPtr(1),
				RestartPolicy: &api.RestartPolicy{
					Attempts: helper.IntToPtr(3),
					Interval: helper.TimeToPtr(5 * time.Minute),
					Delay:    helper.TimeToPtr(25 * time.Second),
					Mode:     helper.StringToPtr("delay"),
				},
				Tasks: []*api.Task{
					&api.Task{
						Name:   beTaskName,
						Driver: "raw_exec",
						Resources: &api.Resources{
							CPU:      helper.IntToPtr(500),
							MemoryMB: helper.IntToPtr(256),
							Networks: []*api.NetworkResource{
								&api.NetworkResource{
									MBits: helper.IntToPtr(400),
								},
							},
						},
						Env: map[string]string{
							"JIVA_REP_NAME":     pvc.Name + "-" + beTaskGroup + "-" + beTaskName,
							"JIVA_CTL_IP":       jivaFeIP,
							"JIVA_REP_VOLNAME":  jivaVolName,
							"JIVA_REP_VOLSIZE":  jivaVolSize,
							"JIVA_REP_VOLSTORE": "/tmp/jiva/" + pvc.Name + beTaskGroup + "/" + beTaskName,
							"JIVA_REP_VERSION":  jivaFeVersion,
							"JIVA_REP_NETWORK":  jivaFeNetwork,
							"JIVA_REP_IFACE":    jivaFeInterface,
							"JIVA_REP_IP":       jivaBeIP,
							"JIVA_REP_SUBNET":   jivaFeSubnet,
						},
						Artifacts: []*api.TaskArtifact{
							&api.TaskArtifact{
								GetterSource: helper.StringToPtr("https://raw.githubusercontent.com/openebs/jiva/master/scripts/launch-jiva-rep-with-ip"),
								RelativeDest: helper.StringToPtr("local/"),
							},
						},
						Config: map[string]interface{}{
							"command": "launch-jiva-rep-with-ip",
						},
						LogConfig: &api.LogConfig{
							MaxFiles:      helper.IntToPtr(3),
							MaxFileSizeMB: helper.IntToPtr(1),
						},
					},
				},
			},
		},
	}, nil
}

// TODO
// Transformation from JobSummary to pv
//
//  1. Need an Interface or functional callback defined at
// lib/api/v1/nomad.go &
//  2. implemented by the volume plugins that want
// to be orchestrated by Nomad
//  3. This transformer instance needs to be injected from
// volume plugin to orchestrator, in a generic way.
func JobSummaryToPv(jobSummary *api.JobSummary) (*v1.PersistentVolume, error) {

	if jobSummary == nil {
		return nil, fmt.Errorf("Nil nomad job summary provided")
	}

	// TODO
	// Needs to be filled up
	return &v1.PersistentVolume{}, nil
}

// TODO
// This transformation is a very crude approach currently.
func JobEvalsToPv(submittedJob *api.Job, evals []*api.Evaluation) (*v1.PersistentVolume, error) {

	if evals == nil {
		return nil, fmt.Errorf("Nil job evaluations provided")
	}

	pvEvals := make([]api.Evaluation, len(evals))

	for i, eval := range evals {
		pvEvals[i] = *eval
	}

	pv := &v1.PersistentVolume{}
	pv.Evals = pvEvals
	pv.Name = *submittedJob.ID

	return pv, nil
}

// Transform a PersistentVolume type to Nomad Job type
func PvToJob(pv *v1.PersistentVolume) (*api.Job, error) {

	if pv == nil {
		return nil, fmt.Errorf("Nil persistent volume provided")
	}

	return &api.Job{
		Name: helper.StringToPtr(pv.Name),
		// TODO
		// ID is same as Name currently
		ID: helper.StringToPtr(pv.Name),
	}, nil
}
