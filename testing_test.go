package resource_test

import (
	"net/http/httptest"
	"testing"
)

var octetMimeType = "application/octet-stream"
var jsonMimeType = "application/json"

func AssertResponseCode(t *testing.T, r *httptest.ResponseRecorder, want int) {
	t.Helper()
	if r.Code != want {
		t.Errorf("got: %v want:%v", r.Code, want)
	}
}

func AssertTrue(t *testing.T, got bool) {
	t.Helper()
	if !got {
		t.Errorf("expecting to be true, got : %v", got)
	}
}

func AssertFalse(t *testing.T, got bool) {
	t.Helper()
	if got {
		t.Errorf("expecting to be false, got : %v", got)
	}
}

func AssertResponseContentType(t *testing.T, response *httptest.ResponseRecorder, mimeType string) {
	t.Helper()
	if response.Header().Get("Content-type") != mimeType {
		t.Errorf("got:%s want:%s", response.Header().Get("Content-type"), mimeType)
	}
}

func AssertNoErrorFatal(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Was not expecting error: %v", err)
	}
}

func AssertError(t *testing.T, err error) {
	if err == nil {
		t.Fatalf("Was expecting error")
	}
}

func AssertEqualError(t *testing.T, err, want error) {
	if err != want {
		t.Errorf("got:%v want:%v", err, want)
	}
}
