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
		panic(&ErrorTypeResourceSlashesNotAllowed{Errorf{FormatErrorResourceSlashesNotAllowed, name}})
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

func (rs *Resource) Get(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodGet, methodOperation, contentTypeSelector)
	return rs.AddMethod(method)
}

func (rs *Resource) Post(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodPost, methodOperation, contentTypeSelector)
	return rs.AddMethod(method)
}

func (rs *Resource) Delete(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodDelete, methodOperation, contentTypeSelector)
	return rs.AddMethod(method)
}

func (rs *Resource) Options(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodOptions, methodOperation, contentTypeSelector)
	return rs.AddMethod(method)
}

func (rs *Resource) Put(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodPut, methodOperation, contentTypeSelector)
	return rs.AddMethod(method)
}

func (rs *Resource) Patch(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodPatch, methodOperation, contentTypeSelector)
	return rs.AddMethod(method)
}

func (rs *Resource) Connect(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodConnect, methodOperation, contentTypeSelector)
	return rs.AddMethod(method)
}

func (rs *Resource) Head(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodHead, methodOperation, contentTypeSelector)
	return rs.AddMethod(method)
}

func (rs *Resource) Trace(methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
	method := NewMethod(http.MethodTrace, methodOperation, contentTypeSelector)
	return rs.AddMethod(method)
}

//Methods returns the collection of methods.
//This is a copy of the internal collection, so methods cannot be changed from this slice
func (rs *Resource) Methods() []Method {
	ms := []Method{}
	for _, m := range rs.methods {
		ms = append(ms, *m)
	}
	return ms
}

// //Resources returns the collection of the resources.
// //This is a copy of the internal collection, so resources cannot be changed from this slice.
// func (rs *Resource) Resources() []Resource {
// 	res := []Resource{}
// 	for _, r := range rs.resources {
// 		res = append(res, r)
// 	}
// 	return res
// }

//Path returns the name and path property.
func (rs *Resource) Path() string {
	return rs.path
}

func (rs *Resource) AddMethod(method Method) *Method {
	if _, ok := rs.methods[method.HTTPMethod]; ok {
		panic(ErrorResourceMethodCollition)
	}
	rs.methods[method.HTTPMethod] = &method
	return rs.methods[method.HTTPMethod]
}

// //Resource creates a new Resource and append resources defined in fn function to the collection of resources to the new resource.
// func (rs *Resource) Resource(name string, fn func(r *Resource)) {
// 	newResource := NewResource(name)
// 	if fn != nil {
// 		fn(&newResource)
// 	}
// 	rs.resources[name] = newResource
// }
