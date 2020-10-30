package resource

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// RestAPI is the root of a REST API abstraction.
// It is responsable for specification generation (GenerateSpec function), and
// Server handler generation (GenerateServer function)
type RestAPI struct {
	ID          string
	Version     string
	Description string
	Title       string
	Host        string
	BasePath    string
	ResourceCollection
}

// NewRestAPI creates a new RestAPI.
// host parameter should not be an URL, but the server host name
func NewRestAPI(basePath, host, title, version string) RestAPI {
	return RestAPI{BasePath: basePath, Host: host, Title: title, Version: version}
}

// GenerateSpec will generate the API specification using APISpecGenerator interface implementation (specGenerator),
// and will write into a io.Writer implementation (writer)
func (r RestAPI) GenerateSpec(writer io.Writer, specGenerator APISpecGenerator) {
	specGenerator.GenerateAPISpec(writer, r)
}

// GenerateServer generates a http.Handler using a ServerGenerator implementation (serverGenerator)
func (r RestAPI) GenerateServer(serverGenerator ServerGenerator) http.Handler {
	resourcesCheck(r.resources)
	server := serverGenerator.GenerateServer(r)
	return inputGetFunctionsMiddleware(serverGenerator.GetURIParam(), server)
}

func inputGetFunctionsMiddleware(getURIParamFunc func(r *http.Request, key string) string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), InputContextKey("uriparamfunc"), getURIParamFunc)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func resourcesCheck(res map[string]Resource) {
	for _, resource := range res {
		for _, m := range resource.methods {
			for _, resp := range m.Responses() {
				httpResponseCodeCheck(resp.Code(), m.HTTPMethod, resource.path)
				parameterOperationCheck(m, resource.path)
			}
		}
		resourcesCheck(resource.resources)
	}
}

// An invalid code will panic in an implementation of http server (see checkWriteHeaderCode function on https://golang.org/src/net/http/server.go)
// We will check this before the server is up and running, and avoid an unexpected panic.
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
