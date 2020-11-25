package rest_test

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"

	"github.com/ehsoc/rest"
	"github.com/ehsoc/rest/encdec"
)

var testNegotiateEncoder = []struct {
	accept            string
	wantedContentType string
}{
	{"application/xml; indent=4, application/json, application/yaml, text/html, */*", "application/xml"},
	{"application/yaml, text/html, application/octet-stream, */*", "application/octet-stream"},
	{"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8", "application/xml"},
	{"text/html, application/xml;q=0.9, application/xhtml+xml, image/png, image/webp, image/jpeg, image/gif, image/x-xbitmap, */*;q=0.1", "application/xml"},
	{"image/jpeg, application/x-ms-application, image/gif, application/xaml+xml, image/pjpeg, application/x-ms-xbap, application/x-shockwave-flash, application/msword, */*", "application/json"},
	{"image/jpeg, application/octet-stream, application/xml ,image/gif", "application/octet-stream"},
	{"", "application/json"},
	{"*/*", "application/json"},
}

func TestNegotiateEncoder(t *testing.T) {
	cts := mustGetCTS()

	for _, tt := range testNegotiateEncoder {
		t.Run("", func(t *testing.T) {
			n := rest.DefaultNegotiator{}
			request, _ := http.NewRequest(http.MethodPost, "/", nil)
			request.Header.Set("Accept", tt.accept)
			got, _, _ := n.NegotiateEncoder(request, &cts)
			if got != tt.wantedContentType {
				t.Errorf("got:%s want:%s", got, tt.wantedContentType)
			}
		})
	}
	t.Run("no default content-type", func(t *testing.T) {
		n := rest.DefaultNegotiator{}
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		ctsnd := mustGetCTSNoDefault()
		_, _, err := n.NegotiateEncoder(request, &ctsnd)
		if err == nil {
			t.Errorf("Was expecting an error")
		}
	})
	t.Run("no content-type registered", func(t *testing.T) {
		n := rest.DefaultNegotiator{}
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		void := rest.NewContentTypes()
		_, _, err := n.NegotiateEncoder(request, &void)
		if err == nil {
			t.Errorf("Was expecting an error")
		}
	})
}

func mustGetCTS() rest.ContentTypes {
	cts := rest.NewContentTypes()
	cts.Add(octetMimeType, &EncodeDecoderSpy{}, false)
	cts.Add("application/json", encdec.JSONEncoderDecoder{}, true)
	cts.Add("application/xml", &EncodeDecoderSpy{}, false)
	return cts
}

func mustGetCTSNoDefault() rest.ContentTypes {
	cts := rest.NewContentTypes()
	cts.Add(octetMimeType, &EncodeDecoderSpy{}, false)
	cts.Add("application/json", encdec.JSONEncoderDecoder{}, false)
	cts.Add("application/xml", &EncodeDecoderSpy{}, false)
	return cts
}

func TestNegotiateDecoder(t *testing.T) {
	cts := mustGetCTS()

	t.Run("with body, no content-type, nodefault content-type", func(t *testing.T) {
		body := bytes.NewBufferString("Not empty body")
		n := rest.DefaultNegotiator{}
		request, _ := http.NewRequest(http.MethodPost, "/", body)
		ctsnd := mustGetCTSNoDefault()
		_, _, err := n.NegotiateDecoder(request, &ctsnd)
		if err == nil {
			t.Error("Was expecting an error")
		}
	})
	t.Run("with body, no content-type", func(t *testing.T) {
		body := bytes.NewBufferString("Not empty body")
		n := rest.DefaultNegotiator{}
		request, _ := http.NewRequest(http.MethodPost, "/", body)
		_, _, err := n.NegotiateDecoder(request, &cts)
		if err == nil {
			t.Error("Was expecting an error")
		}
	})
	t.Run("no body, no content-type", func(t *testing.T) {
		n := rest.DefaultNegotiator{}
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		_, _, err := n.NegotiateDecoder(request, &cts)
		if err != nil {
			t.Fatalf("Was not expecting an error")
		}
	})
	t.Run("with body, with unavailable content-type", func(t *testing.T) {
		body := bytes.NewBufferString("Not empty body")
		n := rest.DefaultNegotiator{}
		request, _ := http.NewRequest(http.MethodPost, "/", body)
		request.Header.Set("Content-type", "unavailable/type")
		_, _, err := n.NegotiateDecoder(request, &cts)
		if err == nil {
			t.Error("Was expecting an error")
		}
	})
	t.Run("with body, with multiple unavailable content-type", func(t *testing.T) {
		body := bytes.NewBufferString("Not empty body")
		n := rest.DefaultNegotiator{}
		request, _ := http.NewRequest(http.MethodPost, "/", body)
		request.Header.Set("Content-type", "unavailable/type, text/xml, text/yaml")
		_, _, err := n.NegotiateDecoder(request, &cts)
		if err == nil {
			t.Error("Was expecting an error")
		}
	})
	t.Run("with body, with blank content-type", func(t *testing.T) {
		n := rest.DefaultNegotiator{}
		body := bytes.NewBufferString("Not empty body")
		request, _ := http.NewRequest(http.MethodPost, "/", body)
		request.Header.Set("Content-type", "")
		_, _, err := n.NegotiateDecoder(request, &cts)
		if err == nil {
			t.Error("Was expecting an error")
		}
	})
	t.Run("with body, with multiple unavailable content-type and one available", func(t *testing.T) {
		body := bytes.NewBufferString("Not empty body")
		n := rest.DefaultNegotiator{}
		request, _ := http.NewRequest(http.MethodPost, "/", body)
		wantMIMETypeName := "application/xml"
		request.Header.Set("Content-type", "unavailable/type, application/xml, text/yaml")
		gotContentTypeName, dec, err := n.NegotiateDecoder(request, &cts)
		if err != nil {
			t.Fatalf("Was not expecting an error")
		}
		if gotContentTypeName != wantMIMETypeName {
			t.Errorf("got:%s want:%s", gotContentTypeName, wantMIMETypeName)
		}
		wantDecoder := &EncodeDecoderSpy{}
		if !reflect.DeepEqual(dec, wantDecoder) {
			t.Errorf("got:%v want%v", dec, wantDecoder)
		}
	})
}
