package resource

import (
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
