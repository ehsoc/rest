package rest

import (
	"net/http"
	"strings"
)

// Resource represents a REST API resource, and a node in a URL tree path.
// It contains a collection of resource methods, and a collection of resources.
type Resource struct {
	path        string
	Summary     string
	Description string
	// a unique method key is defined by a combination of a path and a HTTP method.
	methods map[string]*Method
	ResourceCollection
}

// Middleware is a function that gets a `next` Handler and returns the middleware Handler
type Middleware func(http.Handler) http.Handler

// NewResource creates a new resource node.
// `name` parameter value should not contain any reserved chars like ":/?#[]@!$&'()*+,;=" (RFC 3986 https://tools.ietf.org/html/rfc3986#section-2.2)
// nor curly brackets "{}"
func NewResource(name string) Resource {
	err := validateResourceName(name)
	if err != nil {
		panic(err)
	}
	name = strings.TrimSpace(name)
	r := Resource{}
	r.methods = make(map[string]*Method)
	r.resources = make(map[string]Resource)
	r.path = name
	return r
}

// NewResourceP creates a new URI parameter resource node.
// p Parameter must be URIParameter type. Use NewURIParameter to create one.
func NewResourceP(p Parameter) Resource {
	err := validateResourceName(p.Name)
	if err != nil {
		panic(err)
	}
	r := Resource{}
	r.methods = make(map[string]*Method)
	r.resources = make(map[string]Resource)
	r.path = "{" + strings.TrimSpace(p.Name) + "}"
	return r
}

func validateResourceName(name string) error {
	if strings.ContainsAny(name, ResourceReservedChar) {
		return &TypeErrorResourceCharNotAllowed{errorf{messageErrResourceCharNotAllowed, name}}
	}
	return nil
}

// Get adds a new GET method to the method collection
func (rs *Resource) Get(methodOperation MethodOperation, ct ContentTypes) *Method {
	method := NewMethod(http.MethodGet, methodOperation, ct)
	rs.AddMethod(method)
	return method
}

// Post adds a new POST method to the method collection
func (rs *Resource) Post(methodOperation MethodOperation, ct ContentTypes) *Method {
	method := NewMethod(http.MethodPost, methodOperation, ct)
	rs.AddMethod(method)
	return method
}

// Delete adds a new DELETE method to the method collection
func (rs *Resource) Delete(methodOperation MethodOperation, ct ContentTypes) *Method {
	method := NewMethod(http.MethodDelete, methodOperation, ct)
	rs.AddMethod(method)
	return method
}

// Options adds a new OPTIONS method to the method collection
func (rs *Resource) Options(methodOperation MethodOperation, ct ContentTypes) *Method {
	method := NewMethod(http.MethodOptions, methodOperation, ct)
	rs.AddMethod(method)
	return method
}

// Put adds a new PUT method to the method collection
func (rs *Resource) Put(methodOperation MethodOperation, ct ContentTypes) *Method {
	method := NewMethod(http.MethodPut, methodOperation, ct)
	rs.AddMethod(method)
	return method
}

// Patch adds a new PATCH method to the method collection
func (rs *Resource) Patch(methodOperation MethodOperation, ct ContentTypes) *Method {
	method := NewMethod(http.MethodPatch, methodOperation, ct)
	rs.AddMethod(method)
	return method
}

// Connect adds a new CONNECT method to the method collection
func (rs *Resource) Connect(methodOperation MethodOperation, ct ContentTypes) *Method {
	method := NewMethod(http.MethodConnect, methodOperation, ct)
	rs.AddMethod(method)
	return method
}

// Head adds a new HEAD method to the method collection
func (rs *Resource) Head(methodOperation MethodOperation, ct ContentTypes) *Method {
	method := NewMethod(http.MethodHead, methodOperation, ct)
	rs.AddMethod(method)
	return method
}

// Trace adds a new TRACE method to the method collection
func (rs *Resource) Trace(methodOperation MethodOperation, ct ContentTypes) *Method {
	method := NewMethod(http.MethodTrace, methodOperation, ct)
	rs.AddMethod(method)
	return method
}

// Methods returns the method collection
func (rs *Resource) Methods() []Method {
	rs.checkNilMethods()
	ms := []Method{}
	for _, m := range rs.methods {
		ms = append(ms, *m)
	}

	return ms
}

// Path returns the name and path property.
func (rs *Resource) Path() string {
	return rs.path
}

// AddMethod adds a new method to the method collection.
// If the same HTTPMethod (POST, GET, etc) is already in the collection, will be replaced silently.
// The current resource's middleware stack will be applied
func (rs *Resource) AddMethod(method *Method) {
	rs.checkNilMethods()
	// prepend resource middlewares to themethod
	method.middleware = append(rs.middleware, method.middleware...)
	// replace the core security middleware
	if rs.overWriteCoreSecurityMiddleware != nil {
		method.replaceSecurityMiddleware(rs.overWriteCoreSecurityMiddleware)
	}
	method.buildHandler()
	rs.methods[strings.ToUpper(method.HTTPMethod)] = method
}

func (rs *Resource) checkNilMethods() {
	if rs.methods == nil {
		rs.methods = make(map[string]*Method)
	}
}

// Use adds one or more middlewares to the resources's middleware stack.
// This middleware stack will be applied to the resource methods declared after the call of `Use`.
// The stack will be passed down to the child resources.
func (rs *Resource) Use(m ...Middleware) {
	rs.middleware = append(rs.middleware, m...)
}

// OverwriteCoreSecurityMiddleware will replace the core default security middleware
// of all the child methods and resources declared after the call of this method.
func (rs *Resource) OverwriteCoreSecurityMiddleware(m Middleware) {
	rs.overWriteCoreSecurityMiddleware = m
}
