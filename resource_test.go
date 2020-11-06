package resource_test

import (
	"net/http"
	"reflect"
	"sort"
	"testing"

	"github.com/ehsoc/resource"
)

var moTest = resource.NewMethodOperation(&OperationStub{}, resource.NewResponse(200)).WithFailResponse(resource.NewResponse(404))

func TestAddMethod(t *testing.T) {
	renderers := resource.NewRenderers()
	t.Run("get method", func(t *testing.T) {
		r := resource.NewResource("pet")
		m := resource.NewMethod("GET", moTest, renderers)
		r.AddMethod(m)
		if len(r.Methods()) != 1 {
			t.Errorf("expecting on method")
		}
		getMethod := r.Methods()[0]
		if getMethod.HTTPMethod != http.MethodGet {
			t.Errorf("got: %v want: %v", getMethod.HTTPMethod, http.MethodGet)
		}
		if !reflect.DeepEqual(getMethod.MethodOperation, moTest) {
			t.Errorf("got: %v want: %v", getMethod.MethodOperation, moTest)
		}
	})
	t.Run("methods map nil", func(t *testing.T) {
		defer assertNoPanic(t)
		r := resource.Resource{}
		r.Get(resource.MethodOperation{}, resource.Renderers{})
	})
	t.Run("methods map nil with constructor", func(t *testing.T) {
		defer assertNoPanic(t)
		r := resource.NewResource("pet")
		r.Get(resource.MethodOperation{}, resource.Renderers{})
	})
	t.Run("override existing method", func(t *testing.T) {
		description := "second get"
		r := resource.NewResource("")
		r.Get(resource.MethodOperation{}, resource.Renderers{}).WithDescription("first get")
		r.Post(resource.MethodOperation{}, resource.Renderers{})
		r.Get(resource.MethodOperation{}, resource.Renderers{}).WithDescription(description)
		if len(r.Methods()) != 2 {
			t.Fatalf("expecting 2 methods, got: %v", len(r.Methods()))
		}
		for _, m := range r.Methods() {
			if m.HTTPMethod == "GET" {
				if m.Description != "second get" {
					t.Errorf("got: %v want: %v", m.Description, description)
				}
			}
		}
	})
}

func TestGet(t *testing.T) {
	renderers := resource.NewRenderers()
	t.Run("get method", func(t *testing.T) {
		r := resource.NewResource("pet")
		r.Get(moTest, renderers)
		if len(r.Methods()) != 1 {
			t.Errorf("expecting on method")
		}
		getMethod := r.Methods()[0]
		if getMethod.HTTPMethod != http.MethodGet {
			t.Errorf("got: %v want: %v", getMethod.HTTPMethod, http.MethodGet)
		}
		if !reflect.DeepEqual(getMethod.MethodOperation, moTest) {
			t.Errorf("got: %v want: %v", getMethod.MethodOperation, moTest)
		}
	})
	t.Run("post method", func(t *testing.T) {
		r := resource.NewResource("pet")
		r.Post(moTest, renderers)
		if len(r.Methods()) != 1 {
			t.Errorf("expecting on method")
		}
		getMethod := r.Methods()[0]
		if getMethod.HTTPMethod != http.MethodPost {
			t.Errorf("got: %v want: %v", getMethod.HTTPMethod, http.MethodPost)
		}
		if !reflect.DeepEqual(getMethod.MethodOperation, moTest) {
			t.Errorf("got: %v want: %v", getMethod.MethodOperation, moTest)
		}
	})
	t.Run("options put", func(t *testing.T) {
		r := resource.NewResource("pet")
		r.Put(moTest, renderers)
		if len(r.Methods()) != 1 {
			t.Errorf("expecting on method")
		}
		getMethod := r.Methods()[0]
		if getMethod.HTTPMethod != http.MethodPut {
			t.Errorf("got: %v want: %v", getMethod.HTTPMethod, http.MethodPut)
		}
		if !reflect.DeepEqual(getMethod.MethodOperation, moTest) {
			t.Errorf("got: %v want: %v", getMethod.MethodOperation, moTest)
		}
	})
	t.Run("options patch", func(t *testing.T) {
		r := resource.NewResource("pet")
		r.Patch(moTest, renderers)
		if len(r.Methods()) != 1 {
			t.Errorf("expecting on method")
		}
		getMethod := r.Methods()[0]
		if getMethod.HTTPMethod != http.MethodPatch {
			t.Errorf("got: %v want: %v", getMethod.HTTPMethod, http.MethodPatch)
		}
		if !reflect.DeepEqual(getMethod.MethodOperation, moTest) {
			t.Errorf("got: %v want: %v", getMethod.MethodOperation, moTest)
		}
	})
	t.Run("delete method", func(t *testing.T) {
		r := resource.NewResource("pet")
		r.Delete(moTest, renderers)
		if len(r.Methods()) != 1 {
			t.Errorf("expecting on method")
		}
		getMethod := r.Methods()[0]
		if getMethod.HTTPMethod != http.MethodDelete {
			t.Errorf("got: %v want: %v", getMethod.HTTPMethod, http.MethodDelete)
		}
		if !reflect.DeepEqual(getMethod.MethodOperation, moTest) {
			t.Errorf("got: %v want: %v", getMethod.MethodOperation, moTest)
		}
	})
	t.Run("options method", func(t *testing.T) {
		r := resource.NewResource("pet")
		r.Options(moTest, renderers)
		if len(r.Methods()) != 1 {
			t.Errorf("expecting on method")
		}
		getMethod := r.Methods()[0]
		if getMethod.HTTPMethod != http.MethodOptions {
			t.Errorf("got: %v want: %v", getMethod.HTTPMethod, http.MethodOptions)
		}
		if !reflect.DeepEqual(getMethod.MethodOperation, moTest) {
			t.Errorf("got: %v want: %v", getMethod.MethodOperation, moTest)
		}
	})
	t.Run("connect method", func(t *testing.T) {
		r := resource.NewResource("pet")
		r.Connect(moTest, renderers)
		if len(r.Methods()) != 1 {
			t.Errorf("expecting on method")
		}
		getMethod := r.Methods()[0]
		if getMethod.HTTPMethod != http.MethodConnect {
			t.Errorf("got: %v want: %v", getMethod.HTTPMethod, http.MethodConnect)
		}
		if !reflect.DeepEqual(getMethod.MethodOperation, moTest) {
			t.Errorf("got: %v want: %v", getMethod.MethodOperation, moTest)
		}
	})
	t.Run("head method", func(t *testing.T) {
		r := resource.NewResource("pet")
		r.Head(moTest, renderers)
		if len(r.Methods()) != 1 {
			t.Errorf("expecting on method")
		}
		getMethod := r.Methods()[0]
		if getMethod.HTTPMethod != http.MethodHead {
			t.Errorf("got: %v want: %v", getMethod.HTTPMethod, http.MethodHead)
		}
		if !reflect.DeepEqual(getMethod.MethodOperation, moTest) {
			t.Errorf("got: %v want: %v", getMethod.MethodOperation, moTest)
		}
	})
	t.Run("trace method", func(t *testing.T) {
		r := resource.NewResource("pet")
		r.Trace(moTest, renderers)
		if len(r.Methods()) != 1 {
			t.Errorf("expecting on method")
		}
		getMethod := r.Methods()[0]
		if getMethod.HTTPMethod != http.MethodTrace {
			t.Errorf("got: %v want: %v", getMethod.HTTPMethod, http.MethodTrace)
		}
		if !reflect.DeepEqual(getMethod.MethodOperation, moTest) {
			t.Errorf("got: %v want: %v", getMethod.MethodOperation, moTest)
		}
	})
}

func TestNewResource(t *testing.T) {
	t.Run("new resource", func(t *testing.T) {
		name := "pet"
		resource.NewResource(name)
		defer assertNoPanic(t)
	})
	t.Run("new resource with slash", func(t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				if _, ok := err.(*resource.TypeErrorResourceSlashesNotAllowed); !ok {
					t.Fatalf("got: %T want: %T", err, resource.TypeErrorResourceSlashesNotAllowed{})
				}
			}
		}()
		name := "/pet"
		resource.NewResource(name)
	})
}

func TestResourceIntegration(t *testing.T) {
	r := resource.NewResource("car")
	r.Resource("find", func(r *resource.Resource) {
		r.Resource("left", func(r *resource.Resource) {
		})
		r.Resource("right", func(r *resource.Resource) {
		})
	})
	if r.Path() != "car" {
		t.Errorf("got : %v want: %v", r.Path(), "car")
	}
	findNode := r.Resources()[0]
	if findNode.Path() != "find" {
		t.Errorf("got : %v want: %v", findNode.Path(), "find")
	}
	if len(findNode.Resources()) != 2 {
		t.Fatalf("expecting 2 sub nodes got: %v", len(findNode.Resources()))
	}
	directionResources := findNode.Resources()
	sort.Slice(directionResources, func(i, j int) bool {
		return directionResources[i].Path() < directionResources[j].Path()
	})
	if directionResources[0].Path() != "left" {
		t.Errorf("got : %v want: %v", findNode.Resources()[0].Path(), "left")
	}
	if directionResources[1].Path() != "right" {
		t.Errorf("got : %v want: %v", findNode.Resources()[1].Path(), "right")
	}
}
