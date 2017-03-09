package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/openebs/mayaserver/lib/api/v1"
	"github.com/openebs/mayaserver/lib/volume/jiva"
)

// VolumeSpecificRequest is a http handler implementation.
// The URL path is parsed to match specific implementations.
//
// TODO
//    Should it return specific types than interface{} ?
func (s *HTTPServer) VolumeSpecificRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {

	path := strings.TrimPrefix(req.URL.Path, "/latest/volume")

	// Is req valid ?
	if path == req.URL.Path {
		return nil, CodedError(405, ErrInvalidMethod)
	}

	switch {

	case strings.Contains(path, "/provision/"):
		volName := strings.TrimPrefix(path, "/provision/")
		return s.volumeProvision(resp, req, volName)

	case strings.Contains(path, "/delete/"):
		volName := strings.TrimPrefix(path, "/delete/")
		return s.volumeDelete(resp, req, volName)

	default:
		return nil, CodedError(405, ErrInvalidMethod)
	}
}

func (s *HTTPServer) volumeProvision(resp http.ResponseWriter, req *http.Request, volName string) (interface{}, error) {

	if volName == "" {
		return nil, fmt.Errorf("Volume name missing for provisioning")
	}

	// TODO
	// Get the type of volume plugin from:
	//  1. http request parameters,
	//  2. Mayaconfig etc.
	//
	// We shall hardcode to jiva now
	volPlugName := jiva.JivaStorPluginName

	// Get jiva storage plugin
	jivaStor, err := s.GetVolumePlugin(volPlugName)

	// Get jiva volume provisioner
	jivaProv, ok := jivaStor.Provisioner()
	if !ok {
		return nil, fmt.Errorf("Volume provisioning not supported by '%s'", volPlugName)
	}

	// Provision a jiva volume
	pvc := &v1.PersistentVolumeClaim{}
	pvc.Name = volName

	// TODO
	// This should be set from http query parameters if present
	// Iterate through the query parameters & set them as-is into Labels.
	// Set the empty properties with defaults at respective volume plugin
	// The volume plugin may use a config file to fetch the default values
	//
	// NOTE:
	//  1. The datacenter property should accept multiple values
	//  2. A region can consist of multiple datacenters otherwise known as zones
	pvc.Labels = map[string]string{
		"region":          "global",
		"datacenter":      "dc1",
		"jivafeversion":   "openebs/jiva:latest",
		"jivafenetwork":   "host_static",
		"jivafeip":        "172.28.128.101",
		"jivabeip":        "172.28.128.102",
		"jivafesubnet":    "24",
		"jivafeinterface": "enp0s8",
	}

	pv, err := jivaProv.Provision(pvc)

	if err != nil {
		return nil, err
	}

	return pv, nil
}

func (s *HTTPServer) volumeDelete(resp http.ResponseWriter, req *http.Request, volName string) (interface{}, error) {

	if volName == "" {
		return nil, fmt.Errorf("Volume name missing for deletion")
	}

	// TODO
	// Get the type of volume plugin from:
	//  1. http request parameters,
	//  2. Mayaconfig etc.
	// We shall hardcode to jiva now
	volPlugName := jiva.JivaStorPluginName

	// Get jiva storage plugin
	jivaStor, err := s.GetVolumePlugin(volPlugName)

	// Get jiva volume deleter
	jivaDel, ok := jivaStor.Deleter()
	if !ok {
		return nil, fmt.Errorf("Deleting volume is not supported by '%s'", volPlugName)
	}

	// Delete a jiva volume
	pv := &v1.PersistentVolume{}
	pv.Name = volName

	dPV, err := jivaDel.Delete(pv)

	if err != nil {
		return nil, err
	}

	return dPV, nil
}
