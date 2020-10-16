package resource

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
	Resources
}

// NewResource creates a new resource node.
// name parameter should not contain a slash, because resource represents a unique node and name is the name of the node path
func NewResource(name string) Resource {
	if strings.ContainsAny(name, "/") {
		panic(&TypeErrorResourceSlashesNotAllowed{Errorf{MessageErrResourceSlashesNotAllowed, name}})
	}
	r := Resource{}
	r.methods = make(map[string]*Method)
	r.resources = make(map[string]Resource)
	r.path = name
	return r
}

// Get adds a new GET method to the method collection
func (rs *Resource) Get(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodGet, methodOperation, contentTypeSelector)
	rs.AddMethod(method)
	return method
}

// Post adds a new POST method to the method collection
func (rs *Resource) Post(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodPost, methodOperation, contentTypeSelector)
	rs.AddMethod(method)
	return method
}

// Delete adds a new DELETE method to the method collection
func (rs *Resource) Delete(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodDelete, methodOperation, contentTypeSelector)
	rs.AddMethod(method)
	return method
}

// Options adds a new OPTIONS method to the method collection
func (rs *Resource) Options(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodOptions, methodOperation, contentTypeSelector)
	rs.AddMethod(method)
	return method
}

// Put adds a new PUT method to the method collection
func (rs *Resource) Put(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodPut, methodOperation, contentTypeSelector)
	rs.AddMethod(method)
	return method
}

// Patch adds a new PATCH method to the method collection
func (rs *Resource) Patch(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodPatch, methodOperation, contentTypeSelector)
	rs.AddMethod(method)
	return method
}

// Connect adds a new CONNECT method to the method collection
func (rs *Resource) Connect(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodConnect, methodOperation, contentTypeSelector)
	rs.AddMethod(method)
	return method
}

// Head adds a new HEAD method to the method collection
func (rs *Resource) Head(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodHead, methodOperation, contentTypeSelector)
	rs.AddMethod(method)
	return method
}

// Trace adds a new TRACE method to the method collection
func (rs *Resource) Trace(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodTrace, methodOperation, contentTypeSelector)
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

// AddMethod adds a new method to the method collection
func (rs *Resource) AddMethod(method *Method) {
	rs.checkNilMethods()
	if _, ok := rs.methods[method.HTTPMethod]; ok {
		panic(ErrorResourceMethodCollition)
	}
	rs.methods[method.HTTPMethod] = method
}

func (rs *Resource) checkNilMethods() {
	if rs.methods == nil {
		rs.methods = make(map[string]*Method)
	}
}
