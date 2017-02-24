package server

import (
	"net/http"
	"strings"
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
	//args := structs.VolumeListRequest{}
	//if s.parse(resp, req, &args.Region, &args.QueryOptions) {
	//	return nil, nil
	//}

	//var out structs.VolumeListResponse
	//if err := s.agent.RPC("Job.List", &args, &out); err != nil {
	//	return nil, err
	//}

	//setMeta(resp, &out.QueryMeta)
	//if out.Jobs == nil {
	//	out.Jobs = make([]*structs.JobListStub, 0)
	//}
	//return out.Jobs, nil
	return nil, nil
}

func (s *HTTPServer) volumeDelete(resp http.ResponseWriter, req *http.Request, volName string) (interface{}, error) {
	//args := structs.VolumeListRequest{}
	//if s.parse(resp, req, &args.Region, &args.QueryOptions) {
	//	return nil, nil
	//}

	//var out structs.VolumeListResponse
	//if err := s.agent.RPC("Job.List", &args, &out); err != nil {
	//	return nil, err
	//}

	//setMeta(resp, &out.QueryMeta)
	//if out.Jobs == nil {
	//	out.Jobs = make([]*structs.JobListStub, 0)
	//}
	//return out.Jobs, nil
	return nil, nil
}
