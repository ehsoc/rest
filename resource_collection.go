package rest

// ResourceCollection encapsulate a collection of resource nodes and the methods to add new ones.
// Each node name is unique, in case of conflict the new node will replace the old one silently
type ResourceCollection struct {
	resources map[string]Resource
	// middleware slice is a temporary description of the middleware stack to be applied
	// by a method or other sub-resources
	middleware []Middleware
	// overWriteCoreSecurityMiddleware value nil means default core middleware is applied
	overWriteCoreSecurityMiddleware Middleware
}

// Resources returns the collection of the resource nodes.
// The order of the elements will not be consistent.
func (rs *ResourceCollection) Resources() []Resource {
	rs.checkMap()
	res := []Resource{}

	for _, r := range rs.resources {
		res = append(res, r)
	}
	return res
}

// Resource creates a new resource node and append resources defined in fn function to the collection of resources to the new resource node.
// The usage for the method is as follows:
//
//	r := rest.NewResource("root")
// 	r.Resource("parent", func(r *rest.Resource) {
// 		r.Resource("child1", func(r *rest.Resource) {
// 		})
// 	r.Resource("child2", func(r *rest.Resource) {
// 		})
// 	})
func (rs *ResourceCollection) Resource(name string, fn func(r *Resource)) {
	newResource := NewResource(name)
	rs.addResource(&newResource)
	if fn != nil {
		fn(&newResource)
	}
}

// ResourceP adds a new Resource with a URI parameter path.
// p Parameter must be URIParameter type, use NewURIParameter to create one.
func (rs *ResourceCollection) ResourceP(p Parameter, fn func(r *Resource)) {
	newResource := NewResourceP(p)
	rs.addResource(&newResource)
	if fn != nil {
		fn(&newResource)
	}
}

// addResource adds a resource to the resource collection, the parent resource's middleware stack
// will be prepend to the added resource's stack, also if the resource overWriteCoreSecurityMiddleware is nil
// the parent one will be passed to the child, the child security middleware always override the parent one.
func (rs *ResourceCollection) addResource(r *Resource) {
	// prepend middleware from parent
	r.middleware = append(rs.middleware, r.middleware...)
	// pass the coreSecurityMiddleware if the new resource doesn't have one
	if r.overWriteCoreSecurityMiddleware == nil {
		r.overWriteCoreSecurityMiddleware = rs.overWriteCoreSecurityMiddleware
	}
	rs.checkMap()
	rs.resources[r.path] = *r
}

// checkMap initialize the internal map if is nil
func (rs *ResourceCollection) checkMap() {
	if rs.resources == nil {
		rs.resources = make(map[string]Resource)
	}
}
