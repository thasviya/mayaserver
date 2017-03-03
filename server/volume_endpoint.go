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
// TODO
//    Should it return specific types ?
func (s *HTTPServer) VolumeSpecificRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {

	path := strings.TrimPrefix(req.URL.Path, "/latest/volume")

	// Is req valid ?
	if path == req.URL.Path {
		return nil, CodedError(405, ErrInvalidMethod)
	}

	switch {

	case strings.HasSuffix(path, "/provision"):
		volName := strings.TrimSuffix(path, "/provision")
		return s.volumeProvision(resp, req, volName)

	case strings.HasSuffix(path, "/delete"):
		volName := strings.TrimSuffix(path, "/delete")
		return s.volumeDelete(resp, req, volName)

	default:
		return nil, CodedError(405, ErrInvalidMethod)
	}
}

func (s *HTTPServer) volumeProvision(resp http.ResponseWriter, req *http.Request, volName string) (interface{}, error) {

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
		return nil, fmt.Errorf("Provisioning volume is not supported by '%s'", volPlugName)
	}

	// Provision a jiva volume
	pvc := &v1.PersistentVolumeClaim{}
	pvc.Name = volName

	pv, err := jivaProv.Provision(pvc)

	if err != nil {
		return nil, err
	}

	return pv, nil
}

func (s *HTTPServer) volumeDelete(resp http.ResponseWriter, req *http.Request, volName string) (interface{}, error) {

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

	err = jivaDel.Delete(pv)

	if err != nil {
		return nil, err
	}

	return nil, nil
}
