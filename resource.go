package resource

import (
	"strings"
	"sync"
)

type Resource struct {
	Path        string
	Summary     string
	Description string
	//a unique operation as a combination of a path and an HTTP method is allowed
	methods map[string]Method
	mutex   sync.Mutex
}

func NewResource(path string) Resource {
	r := Resource{}
	r.Path = path
	r.methods = map[string]Method{}
	return r
}

func (r *Resource) AddMethod(HttpMethod string, method Method) error {
	if r.methods == nil {
		r.methods = map[string]Method{}
	}
	_, ok := r.methods[strings.ToUpper(HttpMethod)]
	if ok {
		return ErrorResourceMethodCollition
	}
	r.methods[strings.ToUpper(HttpMethod)] = method
	return nil
}

func (r *Resource) GetMethod(HttpMethod string) (Method, bool) {
	method, ok := r.methods[strings.ToUpper(HttpMethod)]
	return method, ok
}
