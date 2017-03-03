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
		Name: helper.StringToPtr(pvc.Name),
		// TODO
		// ID is same as Name currently
		ID:          helper.StringToPtr(pvc.Name),
		Datacenters: []string{"dc1"},
		Constraints: []*api.Constraint{
			api.NewConstraint("kernel.name", "=", "linux"),
		},
		Meta: map[string]string{
			"JIVA_VOLNAME":           "demo-vsm1-vol1",
			"JIVA_VOLSIZE":           "10g",
			"JIVA_FRONTEND_VERSION":  "openebs/jiva:latest",
			"JIVA_FRONTEND_NETWORK":  "host_static",
			"JIVA_FRONTENDIP":        "172.28.128.101",
			"JIVA_FRONTENDSUBNET":    "24",
			"JIVA_FRONTENDINTERFACE": "enp0s8",
		},
		TaskGroups: []*api.TaskGroup{
			// jiva frontend
			&api.TaskGroup{
				Name: helper.StringToPtr("demo-vsm1-fe"),
				RestartPolicy: &api.RestartPolicy{
					Attempts: helper.IntToPtr(3),
					Interval: helper.TimeToPtr(5 * time.Minute),
					Delay:    helper.TimeToPtr(25 * time.Second),
					Mode:     helper.StringToPtr("delay"),
				},
				Tasks: []*api.Task{
					&api.Task{
						Name:   "fe",
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
							"JIVA_CTL_NAME":    "${NOMAD_JOB_NAME}-${NOMAD_TASK_NAME}",
							"JIVA_CTL_VERSION": "${NOMAD_META_JIVA_FRONTEND_VERSION}",
							"JIVA_CTL_VOLNAME": "${NOMAD_META_JIVA_VOLNAME}",
							"JIVA_CTL_VOLSIZE": "${NOMAD_META_JIVA_VOLSIZE}",
							"JIVA_CTL_IP":      "${NOMAD_META_JIVA_FRONTENDIP}",
							"JIVA_CTL_SUBNET":  "${NOMAD_META_JIVA_FRONTENDSUBNET}",
							"JIVA_CTL_IFACE":   "${NOMAD_META_JIVA_FRONTENDINTERFACE}",
						},
						Artifacts: []*api.TaskArtifact{
							&api.TaskArtifact{
								GetterSource: helper.StringToPtr("https://raw.githubusercontent.com/openebs/jiva/master/scripts/launch-jiva-ctl-with-ip"),
							},
						},
						Config: map[string]interface{}{
							"command": "launch-jiva-ctl-with-ip",
						},
					},
				},
			},
			// jiva replica
			&api.TaskGroup{
				Name: helper.StringToPtr("demo-vsm1-backend-container1"),
				RestartPolicy: &api.RestartPolicy{
					Attempts: helper.IntToPtr(3),
					Interval: helper.TimeToPtr(5 * time.Minute),
					Delay:    helper.TimeToPtr(25 * time.Second),
					Mode:     helper.StringToPtr("delay"),
				},
				Tasks: []*api.Task{
					&api.Task{
						Name:   "be-store1",
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
							"JIVA_REP_NAME":     "${NOMAD_JOB_NAME}-${NOMAD_TASK_NAME}",
							"JIVA_CTL_IP":       "${NOMAD_META_JIVA_FRONTENDIP}",
							"JIVA_REP_VOLNAME":  "${NOMAD_META_JIVA_VOLNAME}",
							"JIVA_REP_VOLSIZE":  "${NOMAD_META_JIVA_VOLSIZE}",
							"JIVA_REP_VOLSTORE": "/tmp/jiva/vsm1/rep1",
							"JIVA_REP_VERSION":  "openebs/jiva:latest",
							"JIVA_REP_NETWORK":  "host_static",
							"JIVA_REP_IFACE":    "enp0s8",
							"JIVA_REP_IP":       "172.28.128.102",
							"JIVA_REP_SUBNET":   "24",
						},
						Artifacts: []*api.TaskArtifact{
							&api.TaskArtifact{
								GetterSource: helper.StringToPtr("https://raw.githubusercontent.com/openebs/jiva/master/scripts/launch-jiva-rep-with-ip"),
							},
						},
						Config: map[string]interface{}{
							"command": "launch-jiva-rep-with-ip",
						},
					},
				},
			},
		},
	}, nil
}

func JobSummaryToPv(jobSummary *api.JobSummary) (*v1.PersistentVolume, error) {

	if jobSummary == nil {
		return nil, fmt.Errorf("Nil nomad job summary provided")
	}

	return &v1.PersistentVolume{}, nil
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
