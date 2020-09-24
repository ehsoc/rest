package resource

import (
	"fmt"
	"io"
	"net/http"
)

//RestAPI is the root of a REST API abstraction.
//It is responsable document generation like output Open API v2 json generation and
//Server handler generation
type RestAPI struct {
	ID       string
	Version  string
	Host     string
	BasePath string
	Resources
}

func (r RestAPI) GenerateSpec(w io.Writer, docGenerator APISpecGenerator) {
	docGenerator.GenerateAPISpec(w, r)
}

func (r RestAPI) GenerateServer(d ServerGenerator) http.Handler {
	resourcesCheck(r.resources)
	return d.GenerateServer(r)
}

func resourcesCheck(res map[string]Resource) {
	for _, resource := range res {
		for _, m := range resource.methods {
			for _, resp := range m.Responses {
				httpResponseCodeCheck(resp.Code, m.HTTPMethod, resource.path)
				parameterOperationCheck(m, resource.path)
			}
		}
		resourcesCheck(resource.resources)
	}
}

//An invalid code will panic in an implementation of http server (see checkWriteHeaderCode function on https://golang.org/src/net/http/server.go)
//We will enforce this before the server is up and running, and avoid an unexpected panic.
func httpResponseCodeCheck(code int, httpMethod string, path string) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("GenerateServer check error: invalid response code %v on method: %v of resource: %v", code, httpMethod, path))
	}
}

func parameterOperationCheck(m *Method, path string) {
	if m.MethodOperation.Operation == nil {
		panic(fmt.Sprintf("GenerateServer check error: resource %s method %s doesn't have an operation.", path, m.HTTPMethod))
	}
}
