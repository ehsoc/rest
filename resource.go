package resource

import (
	"encoding/json"
	"strings"
	"sync"
)

type Resource struct {
	Path string
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

func (r *Resource) AddMethod(method string, m Method) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	_, ok := r.methods[method]
	if ok {
		return ErrorResourceMethodCollition
	}
	r.methods[strings.ToLower(method)] = m
	return nil
}

func (r *Resource) GenerateOpenApiJSON() string {
	j, _ := json.Marshal(r)
	return string(j)
}
