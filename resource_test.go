package resource_test

import (
	"net/http"
	"reflect"
	"sort"
	"testing"

	"github.com/ehsoc/resource"
)

const summary string = "this is my method"

//TODO: test every method
func TestGet(t *testing.T) {
	ct := resource.NewHTTPContentTypeSelector(resource.DefaultUnsupportedMediaResponse)
	t.Run("get method", func(t *testing.T) {
		r := resource.NewResource("pet")
		mo := resource.NewMethodOperation(
			&OperationStub{},
			resource.Response{200, nil, ""},
			resource.Response{404, nil, ""},
			true)
		r.Get(mo, ct)
		if len(r.Methods()) != 1 {
			t.Errorf("expecting on method")
		}
		getMethod := r.Methods()[0]
		if getMethod.HTTPMethod != http.MethodGet {
			t.Errorf("got: %v want: %v", getMethod.HTTPMethod, http.MethodGet)
		}
		if !reflect.DeepEqual(getMethod.MethodOperation, mo) {
			t.Errorf("got: %v want: %v", getMethod.MethodOperation, mo)
		}
	})
	t.Run("methods map nil", func(t *testing.T) {
		defer assertNoPanic(t)
		r := resource.NewResource("pet")
		r.Get(resource.MethodOperation{}, resource.HTTPContentTypeSelector{})
	})
	t.Run("error on existing method", func(t *testing.T) {
		defer assertPanicError(t, resource.ErrorResourceMethodCollition)
		r := resource.NewResource("")
		r.Get(resource.MethodOperation{}, resource.HTTPContentTypeSelector{})
		r.Post(resource.MethodOperation{}, resource.HTTPContentTypeSelector{})
		r.Get(resource.MethodOperation{}, resource.HTTPContentTypeSelector{})
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
				if _, ok := err.(resource.ErrorTypeResourceSlashesNotAllowed); !ok {
					t.Fatalf("got: %T want: %T", err, resource.ErrorTypeResourceSlashesNotAllowed{})
				}
			}
		}()
		name := "/pet"
		resource.NewResource(name)
	})
}

func TestResource(t *testing.T) {
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
