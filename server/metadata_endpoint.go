package server

import (
	"net/http"
	"strings"
)

func (s *HTTPServer) MetaSpecificRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	path := strings.TrimPrefix(req.URL.Path, "/latest/meta-data")

	// Is req valid ?
	if path == req.URL.Path {
		return nil, CodedError(405, ErrInvalidMethod)
	}

	switch {
	case strings.HasSuffix(path, "/instance-id"):
		return s.metaInstanceID(resp, req)
	case strings.HasSuffix(path, "/placement/availability-zone"):
		return s.metaAvailabilityZone(resp, req)
	default:
		return nil, CodedError(405, ErrInvalidMethod)
	}
}

// EBS demands a particular instance id to be returned during
// aws session creation.
func (s *HTTPServer) metaInstanceID(resp http.ResponseWriter, req *http.Request) (interface{}, error) {

	if req.Method != "GET" {
		return nil, CodedError(405, ErrInvalidMethod)
	}

	// OpenEBS can be used as a persistence mechanism for
	// any type of compute instance
	out := "any-compute"

	return out, nil
}

func (s *HTTPServer) metaAvailabilityZone(resp http.ResponseWriter, req *http.Request) (interface{}, error) {

	if req.Method != "GET" {
		return nil, CodedError(405, ErrInvalidMethod)
	}

	// TODO We shall see how to construct an Availability Zone
	out := "any-zone"

	return out, nil
}
