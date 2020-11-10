package rest_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ehsoc/rest"
)

type GenStub struct {
	called    bool
	getURIVal string
}

func (g *GenStub) GenerateAPISpec(w io.Writer, r rest.API) {
	g.called = true
}

func (g *GenStub) GenerateServer(r rest.API) http.Handler {
	g.called = true
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getURIValue := r.Context().Value(rest.InputContextKey("uriparamfunc"))
		if getURIParamFunc, ok := getURIValue.(func(r *http.Request, key string) string); ok {
			g.getURIVal = getURIParamFunc(r, "")
		}
	})
}

func (g *GenStub) GetURIParam() func(*http.Request, string) string {
	return func(r *http.Request, p string) string {
		return "my val"
	}
}

func TestGenerateServer(t *testing.T) {
	g := &GenStub{}
	r := rest.API{}
	server := r.GenerateServer(g)
	if !g.called {
		t.Errorf("Expecting function called")
	}
	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	server.ServeHTTP(resp, req)
	assertStringEqual(t, g.getURIVal, "my val")
}

func TestGenerateSpec(t *testing.T) {
	g := &GenStub{}
	r := rest.API{}
	r.GenerateSpec(new(bytes.Buffer), g)
	if !g.called {
		t.Errorf("Expecting function called")
	}
}
