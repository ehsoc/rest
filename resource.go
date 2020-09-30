package resource

import (
	"net/http"
	"regexp"
	"strings"
)

//Resource represents a node in a url path
type Resource struct {
	path        string
	Summary     string
	Description string
	//a unique operation as a combination of a path and an HTTP method is allowed
	methods map[string]*Method
	Resources
	uRIParam Parameter
}

//NewResource creates a new resource node.
//name parameter should not contain a slash, because resource represents a unique node and name is the name of the node path
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

// func newResourceWithURIParam(pathStr string, paramDescription string, paramType reflect.Kind) (Resource, error) {
// 	params := getURLParamName(pathStr)
// 	if params == nil {
// 		return Resource{}, ErrorResourceURIParamNoParamFound
// 	}
// 	if len(params) > 1 {
// 		return Resource{}, ErrorResourceURIParamMoreThanOne
// 	}
// 	r := Resource{}
// 	r.path = pathStr
// 	r.uRIParam = NewURIParameter(strings.Trim(params[0], "{}"), paramType).WithDescription(paramDescription)
// 	return r, nil
// }

func getURLParamName(pathStr string) []string {
	re := regexp.MustCompile(`\{(.*?)\}`)
	return re.FindAllString(pathStr, -1)
}

//Get adds a new method GET to the method collection
func (rs *Resource) Get(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodGet, methodOperation, contentTypeSelector)
	rs.AddMethod(method)
	return method
}

//Post adds a new method POST to the method collection
func (rs *Resource) Post(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodPost, methodOperation, contentTypeSelector)
	rs.AddMethod(method)
	return method
}

//Delete adds a new method DELETE to the method collection
func (rs *Resource) Delete(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodDelete, methodOperation, contentTypeSelector)
	rs.AddMethod(method)
	return method
}

//Options adds a new method OPTIONS to the method collection
func (rs *Resource) Options(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodOptions, methodOperation, contentTypeSelector)
	rs.AddMethod(method)
	return method
}

//Put adds a new method PUT to the method collection
func (rs *Resource) Put(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodPut, methodOperation, contentTypeSelector)
	rs.AddMethod(method)
	return method
}

//Patch adds a new method PATCH to the method collection
func (rs *Resource) Patch(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodPatch, methodOperation, contentTypeSelector)
	rs.AddMethod(method)
	return method
}

//Connect adds a new method CONNECT to the method collection
func (rs *Resource) Connect(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodConnect, methodOperation, contentTypeSelector)
	rs.AddMethod(method)
	return method
}

//Head adds a new method HEAD to the method collection
func (rs *Resource) Head(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodHead, methodOperation, contentTypeSelector)
	rs.AddMethod(method)
	return method
}

//Trace adds a new method TRACE to the method collection
func (rs *Resource) Trace(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodTrace, methodOperation, contentTypeSelector)
	rs.AddMethod(method)
	return method
}

//Methods returns the collection of methods.
//This is a copy of the internal collection, so methods cannot be changed from this slice
func (rs *Resource) Methods() []Method {
	rs.checkNilMethods()
	ms := []Method{}
	for _, m := range rs.methods {
		ms = append(ms, *m)
	}
	return ms
}

//Path returns the name and path property.
func (rs *Resource) Path() string {
	return rs.path
}

//AddMethod adds a new method to the method collection
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
