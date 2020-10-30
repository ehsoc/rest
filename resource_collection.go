package resource

// ResourceCollection encapsulate a collection of resource nodes and the methods to add new ones.
// Each node name is unique, in case of conflict the new node will replace the old one silently
type ResourceCollection struct {
	resources map[string]Resource
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
//	r := resource.NewResource("root")
// 	r.Resource("parent", func(r *resource.Resource) {
// 		r.Resource("child1", func(r *resource.Resource) {
//
// 		})
// 	r.Resource("child2", func(r *resource.Resource) {
//
// 		})
// 	})
func (rs *ResourceCollection) Resource(name string, fn func(r *Resource)) {
	rs.checkMap()
	newResource := NewResource(name)
	if fn != nil {
		fn(&newResource)
	}
	rs.resources[name] = newResource
}

// AddResource adds a resource node to the collection of resources.
// If other node with the same name is in the collection, it will replace it silently
func (rs *ResourceCollection) AddResource(resource Resource) {
	rs.checkMap()
	rs.resources[resource.path] = resource
}

// checkMap initialize the internal map if is nil
func (rs *ResourceCollection) checkMap() {
	if rs.resources == nil {
		rs.resources = make(map[string]Resource)
	}
}
