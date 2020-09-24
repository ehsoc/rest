package resource_test

import (
	"net/http/httptest"
	"testing"
)

var octetMimeType = "application/octet-stream"
var jsonMimeType = "application/json"

func asserStringEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got: %v want: %v", got, want)
	}
}

func assertResponseCode(t *testing.T, r *httptest.ResponseRecorder, want int) {
	t.Helper()
	if r.Code != want {
		t.Errorf("got: %v want: %v", r.Code, want)
	}
}

func assertTrue(t *testing.T, got bool) {
	t.Helper()
	if !got {
		t.Errorf("expecting to be true, got: %v", got)
	}
}

func assertFalse(t *testing.T, got bool) {
	t.Helper()
	if got {
		t.Errorf("expecting to be false, got: %v", got)
	}
}

func assertResponseContentType(t *testing.T, response *httptest.ResponseRecorder, mimeType string) {
	t.Helper()
	if response.Header().Get("Content-type") != mimeType {
		t.Errorf("got:%s want:%s", response.Header().Get("Content-type"), mimeType)
	}
}

func assertNoErrorFatal(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Was not expecting error: %v", err)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("Was expecting error")
	}
}

func assertEqualError(t *testing.T, err, want error) {
	t.Helper()
	if err != want {
		t.Errorf("got:%v want:%v", err, want)
	}
}

func assertPanicError(t *testing.T, want error) {
	t.Helper()
	if err := recover(); err != want {
		t.Fatalf("got: %v want: %v", err, want)
	}
}

func assertNoPanic(t *testing.T) {
	t.Helper()
	if err := recover(); err != nil {
		t.Fatalf("Not expecting panic: %v", err)
	}
}
