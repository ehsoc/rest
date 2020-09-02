package resource_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/ehsoc/resource"
)

type GenStub struct {
	called bool
}

func (g *GenStub) GenerateAPISpec(w io.Writer, r resource.RestAPI) {
	g.called = true
}

func TestGenerateAPISpec(t *testing.T) {
	g := &GenStub{}
	r := resource.RestAPI{}
	r.GenerateAPISpec(new(bytes.Buffer), g)
	if !g.called {
		t.Errorf("Expecting function called")
	}
}
