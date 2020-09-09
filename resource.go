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
	methods   map[string]Method
	Resources []Resource
	uRIParam  *Parameter
}

func NewResource(pathStr string) (Resource, error) {
	if strings.ContainsAny(pathStr, "{}") {
		return Resource{}, ErrorResourceBracketsNotAllowed
	}
	r := Resource{}
	r.Path = pathStr
	r.methods = map[string]Method{}
	return r, nil
}

func NewResourceWithURIParam(pathStr string, getURIParamFunc GetParamFunc, paramDescription string, paramType reflect.Kind) (Resource, error) {
	params := getURLParamName(pathStr)
	if params == nil {
		return Resource{}, ErrorResourceURIParamNoParamFound
	}
	if len(params) > 1 {
		return Resource{}, ErrorResourceURIParamMoreThanOne
	}
	r := Resource{}
	r.Path = pathStr
	r.methods = make(map[string]Method)
	r.uRIParam = NewURIParameter(strings.Trim(params[0], "{}"), paramType, getURIParamFunc).WithDescription(paramDescription)
	return r, nil
}

func (r *Resource) GetURIParam() *Parameter {
	return r.uRIParam
}

func getURLParamName(pathStr string) []string {
	re := regexp.MustCompile(`\{(.*?)\}`)
	return re.FindAllString(pathStr, -1)
}

func (r *Resource) AddMethod(method Method) error {
	if r.methods == nil {
		r.methods = map[string]Method{}
	}
	HTTPMethod := strings.ToUpper(method.HTTPMethod)
	_, ok := r.methods[HTTPMethod]
	if ok {
		return ErrorResourceMethodCollition
	}
	r.methods[HTTPMethod] = method
	return nil
}

// func (r *Resource) AddURIParamResource(path string, gFunc GetParamFunc) (Resource, error) {
// 	newResource := Resource{}
// 	r.Resources = append(r.Resources, newResource)
// 	return newResource, nil
// }

func (r *Resource) GetMethod(HttpMethod string) (Method, bool) {
	method, ok := r.methods[strings.ToUpper(HttpMethod)]
	return method, ok
}
