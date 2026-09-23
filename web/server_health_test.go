package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/injoyai/tdx"
)

func TestHandleHealthCheckRequiresInitializedCodes(t *testing.T) {
	originalCodes := tdx.DefaultCodes
	t.Cleanup(func() {
		tdx.DefaultCodes = originalCodes
	})

	tdx.DefaultCodes = nil
	recorder := httptest.NewRecorder()
	handleHealthCheck(recorder, httptest.NewRequest(http.MethodGet, "/api/health", nil))
	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("got status %d, want %d", recorder.Code, http.StatusServiceUnavailable)
	}

	tdx.DefaultCodes = &tdx.Codes{Map: map[string]*tdx.CodeModel{
		"sh501046": {Code: "501046", Exchange: "sh"},
	}}
	recorder = httptest.NewRecorder()
	handleHealthCheck(recorder, httptest.NewRequest(http.MethodGet, "/api/health", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", recorder.Code, http.StatusOK)
	}
}

func TestInitializedServerReusesManagerCodes(t *testing.T) {
	if manager == nil || manager.Codes == nil {
		t.Fatal("manager code cache is not initialized")
	}
	if tdx.DefaultCodes != manager.Codes {
		t.Fatal("server and manager use different code caches")
	}
}
