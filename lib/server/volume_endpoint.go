package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/openebs/mayaserver/lib/api/v1"
	"github.com/openebs/mayaserver/lib/volume/jiva"
)

func (s *HTTPServer) VolumesRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	switch req.Method {
	case "GET":
		return s.volumeListRequest(resp, req)
	case "PUT", "POST":
		return s.volumeUpdate(resp, req, "")
	default:
		return nil, CodedError(405, ErrInvalidMethod)
	}
}

// TODO
// Not yet implemented
func (s *HTTPServer) volumeListRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	return nil, CodedError(405, "Volume list not yet implemented")
}

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

	case strings.Contains(path, "/delete/"):
		volName := strings.TrimPrefix(path, "/delete/")
		return s.volumeDelete(resp, req, volName)

	default:
		return nil, CodedError(405, ErrInvalidMethod)
	}
}

func (s *HTTPServer) volumeUpdate(resp http.ResponseWriter, req *http.Request, volName string) (interface{}, error) {

	pvc := v1.PersistentVolumeClaim{}

	if err := decodeYamlBody(req, &pvc); err != nil {
		return nil, CodedError(400, err.Error())
	}

	//if pvc == nil {
	//	return nil, CodedError(400, "Empty or Invalid volume claim request")
	//}

	if pvc.Name == "" {
		return nil, CodedError(400, fmt.Sprintf("Volume name hasn't been provided: '%v'", pvc))
	}

	if pvc.Labels == nil {
		return nil, CodedError(400, fmt.Sprintf("Volume labels hasn't been provided: '%v'", pvc))
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

	pv, err := jivaProv.Provision(&pvc)

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
