package rest_test

import (
	"net/http"
	"reflect"
	"sort"
	"testing"

	"github.com/ehsoc/rest"
)

var moTest = rest.NewMethodOperation(&OperationStub{}, rest.NewResponse(200)).WithFailResponse(rest.NewResponse(404))

func TestAddMethod(t *testing.T) {
	ct := rest.NewContentTypes()
	t.Run("get method", func(t *testing.T) {
		r := rest.NewResource("pet")
		m := rest.NewMethod("GET", moTest, ct)
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
		r := rest.Resource{}
		r.Get(rest.MethodOperation{}, rest.ContentTypes{})
	})
	t.Run("methods map nil with constructor", func(t *testing.T) {
		defer assertNoPanic(t)
		r := rest.NewResource("pet")
		r.Get(rest.MethodOperation{}, rest.ContentTypes{})
	})
	t.Run("override existing method", func(t *testing.T) {
		description := "second get"
		r := rest.NewResource("")
		r.Get(rest.MethodOperation{}, rest.ContentTypes{}).WithDescription("first get")
		r.Post(rest.MethodOperation{}, rest.ContentTypes{})
		r.Get(rest.MethodOperation{}, rest.ContentTypes{}).WithDescription(description)
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
	ct := rest.NewContentTypes()
	t.Run("get method", func(t *testing.T) {
		r := rest.NewResource("pet")
		r.Get(moTest, ct)
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
		r := rest.NewResource("pet")
		r.Post(moTest, ct)
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
		r := rest.NewResource("pet")
		r.Put(moTest, ct)
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
		r := rest.NewResource("pet")
		r.Patch(moTest, ct)
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
		r := rest.NewResource("pet")
		r.Delete(moTest, ct)
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
		r := rest.NewResource("pet")
		r.Options(moTest, ct)
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
		r := rest.NewResource("pet")
		r.Connect(moTest, ct)
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
		r := rest.NewResource("pet")
		r.Head(moTest, ct)
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
		r := rest.NewResource("pet")
		r.Trace(moTest, ct)
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
		rest.NewResource(name)
		defer assertNoPanic(t)
	})
	t.Run("new resource with slash", func(t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				if _, ok := err.(*rest.TypeErrorResourceSlashesNotAllowed); !ok {
					t.Fatalf("got: %T want: %T", err, rest.TypeErrorResourceSlashesNotAllowed{})
				}
			}
		}()
		name := "/pet"
		rest.NewResource(name)
	})
}

func TestResourceIntegration(t *testing.T) {
	r := rest.NewResource("car")
	r.Resource("find", func(r *rest.Resource) {
		r.Resource("left", func(r *rest.Resource) {
		})
		r.Resource("right", func(r *rest.Resource) {
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
