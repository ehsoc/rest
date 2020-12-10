package rest

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// API is the root of a REST API abstraction.
// The two main responsibilities are generate specification generation (GenerateSpec function), and
// Server handler generation (GenerateServer function).
type API struct {
	ID          string
	Version     string
	Description string
	Title       string
	Host        string
	BasePath    string
	ResourceCollection
}

const (
	// URIReservedChar is the reserved char set as in RFC 3986 (https://tools.ietf.org/html/rfc3986#section-2.2)
	URIReservedChar = ":/?#[]@!$&'()*+,;="
	// ResourceReservedChar is URIReservedChar + "{}" that represents the charset that is not allowed on resource names.
	// The curly brackets are reserved for internally denoting a path URI parameter, also used by API specification formats.
	ResourceReservedChar = URIReservedChar + "{}"
)

// NewAPI creates a new API.
func NewAPI(basePath, hostname, title, version string) API {
	return API{BasePath: basePath, Host: hostname, Title: title, Version: version}
}

// GenerateSpec will generate the API specification using APISpecGenerator interface implementation (g),
// and will write into a io.Writer implementation (w)
func (a API) GenerateSpec(w io.Writer, g APISpecGenerator) {
	g.GenerateAPISpec(w, a)
}

// GenerateServer generates a http.Handler using a ServerGenerator implementation (g)
func (a API) GenerateServer(g ServerGenerator) http.Handler {
	resourcesCheck(a.resources)
	server := g.GenerateServer(a)

	return inputGetFunctionsMiddleware(g.GetURIParam(), server)
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
