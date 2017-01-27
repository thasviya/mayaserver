package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInvalidMetaDataRequest(t *testing.T) {
	s := makeHTTPTestServer(t, nil)
	defer s.Cleanup()

	resp := httptest.NewRecorder()
	// valid uri is meta-data & not metadata
	req, _ := http.NewRequest("GET", "/metadata/instance-id", nil)

	out, err := s.Server.MetaSpecificRequest(resp, req)

	if err == nil {
		t.Fatalf("ERR: expected: %v, got: %v", ErrInvalidMethod, err)
	}

	if err.Error() != ErrInvalidMethod {
		t.Fatalf("ERR: expected: %v, got: %v", ErrInvalidMethod, err.Error())
	}

	if out != nil {
		t.Fatalf("Service must not return any value, for invalid request")
	}
}

func TestMetaInstanceID(t *testing.T) {
	s := makeHTTPTestServer(t, nil)
	defer s.Cleanup()

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/latest/meta-data/instance-id", nil)

	out, err := s.Server.MetaSpecificRequest(resp, req)

	if err != nil {
		t.Fatalf("ERR: %v", err)
	}

	if out == "" || out == nil {
		t.Fatalf("Service must return a non empty instance")
	}
}

func TestInvalidMetaReqInstanceID(t *testing.T) {
	s := makeHTTPTestServer(t, nil)
	defer s.Cleanup()

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/latest/meta-data/instance-id", nil)

	out, err := s.Server.MetaSpecificRequest(resp, req)

	if err == nil {
		t.Fatalf("ERR: expected: %v, got: %v", CodedError(405, ErrInvalidMethod), err)
	}

	if err.Error() != ErrInvalidMethod {
		t.Fatalf("ERR: expected: %v, got: %v", ErrInvalidMethod, err.Error())
	}

	if out != nil {
		t.Fatalf("Service must not return any value, for invalid request")
	}
}

func TestMetaAvailZone(t *testing.T) {
	s := makeHTTPTestServer(t, nil)
	defer s.Cleanup()

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/latest/meta-data/placement/availability-zone", nil)

	out, err := s.Server.MetaSpecificRequest(resp, req)

	if err != nil {
		t.Fatalf("ERR: %v", err)
	}

	if out == "" || out == nil {
		t.Fatalf("Service must return a non empty instance")
	}
}

func TestInvalidMetaReqAvailZone(t *testing.T) {
	s := makeHTTPTestServer(t, nil)
	defer s.Cleanup()

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/latest/meta-data/placement/availability-zone", nil)

	out, err := s.Server.MetaSpecificRequest(resp, req)

	if err == nil {
		t.Fatalf("ERR: expected: %v, got: %v", CodedError(405, ErrInvalidMethod), err)
	}

	if err.Error() != ErrInvalidMethod {
		t.Fatalf("ERR: expected: %v, got: %v", ErrInvalidMethod, err.Error())
	}

	if out != nil {
		t.Fatalf("Service must not return any value, for invalid request")
	}
}
