package nomad

import (
	"fmt"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/helper"
	"github.com/openebs/mayaserver/lib/api/v1"
)

func PvcToJob(pvc *v1.PersistentVolumeClaim) (*api.Job, error) {

	if pvc == nil {
		return nil, fmt.Errorf("Nil persistent volume claim provided")
	}

	return &api.Job{
		Name: helper.StringToPtr(pvc.Name),
	}, nil
}

func JobSummaryToPv(jobSummary *api.JobSummary) (*v1.PersistentVolume, error) {

	if jobSummary == nil {
		return nil, fmt.Errorf("Nil nomad job summary provided")
	}

	return &v1.PersistentVolume{}, nil
}

func PvToJob(pv *v1.PersistentVolume) (*api.Job, error) {

	if pv == nil {
		return nil, fmt.Errorf("Nil persistent volume provided")
	}

	return &api.Job{}, nil
}
