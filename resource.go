package resource

import (
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

type GetParamFunc func(r *http.Request) string

type Resource struct {
	Path        string
	Summary     string
	Description string
	//a unique operation as a combination of a path and an HTTP method is allowed
	Methods   []Method
	Resources []*Resource
	uRIParam  *Parameter
}

func NewResource(pathStr string) (Resource, error) {
	if strings.ContainsAny(pathStr, "{}") {
		return Resource{}, ErrorResourceBracketsNotAllowed
	}
	r := Resource{}
	r.Path = pathStr
	return r, nil
}

func NewResourceWithURIParam(pathStr string, paramDescription string, paramType reflect.Kind) (Resource, error) {
	params := getURLParamName(pathStr)
	if params == nil {
		return Resource{}, ErrorResourceURIParamNoParamFound
	}
	if len(params) > 1 {
		return Resource{}, ErrorResourceURIParamMoreThanOne
	}
	r := Resource{}
	r.Path = pathStr
	r.uRIParam = NewURIParameter(strings.Trim(params[0], "{}"), paramType).WithDescription(paramDescription)
	return r, nil
}

func (rs *Resource) GetURIParam() *Parameter {
	return rs.uRIParam
}

func getURLParamName(pathStr string) []string {
	re := regexp.MustCompile(`\{(.*?)\}`)
	return re.FindAllString(pathStr, -1)
}

func (rs *Resource) AddMethod(method Method) error {
	if _, ok := rs.GetMethod(method.HTTPMethod); ok {
		return ErrorResourceMethodCollition
	}
	rs.Methods = append(rs.Methods, method)
	return nil
}

func (rs *Resource) GetMethod(HttpMethod string) (Method, bool) {
	for _, m := range rs.Methods {
		if strings.ToUpper(m.HTTPMethod) == strings.ToUpper(HttpMethod) {
			return m, true
		}
	}
	return Method{}, false
}

//Resource creates a new Resource and append resources defined in fn function
func (rs *Resource) Resource(resourcePath string, fn func(r *Resource)) {
	newResource, _ := NewResource(resourcePath)
	if fn != nil {
		fn(&newResource)
	}
	rs.Resources = append(rs.Resources, &newResource)
}

// //Method creates a new method and append it to the Methods property of Resource
// func (rs *Resource) Method(HTTPMethod string, methodOperation MethodOperation, contentTypeSelector HTTPContentTypeSelector) *Method {
// 	newMethod := NewMethod(HTTPMethod, methodOperation, contentTypeSelector)
// 	err := rs.AddMethod(newMethod)
// 	if err != nil {
// 		log.Panicf("resource: %v", err)
// 	}
// 	return &newMethod
// }
