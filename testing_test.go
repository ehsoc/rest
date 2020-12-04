package rest_test

import (
	"net/http/httptest"
	"testing"

	"github.com/ehsoc/rest"
	"github.com/ehsoc/rest/encdec"
)

var octetMimeType = "application/octet-stream"

func assertStringEqual(t *testing.T, got, want string) {
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

func assertNoErrorFatal(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Was not expecting error: %v", err)
	}
}

func assertEqualError(t *testing.T, err, want error) {
	t.Helper()
	if err != want {
		t.Errorf("got:%v want:%v", err, want)
	}
}

func assertNoPanic(t *testing.T) {
	t.Helper()
	if err := recover(); err != nil {
		t.Fatalf("Not expecting panic: %v", err)
	}
}

func mustGetJSONContentType() rest.ContentTypes {
	jsonContentType := rest.NewContentTypes()
	jsonContentType.Add("application/json", encdec.JSONEncoderDecoder{}, true)
	return jsonContentType
}
