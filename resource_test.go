package resource_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/ehsoc/resource"
)

const summary string = "this is my method"

func TestAddMethod(t *testing.T) {
	t.Run("methods map nil", func(t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Not expecting panic: %v", err)
			}
		}()
		r := resource.Resource{}
		r.AddMethod("POST", resource.Method{})
	})
	t.Run("error on existing key", func(t *testing.T) {
		r := resource.Resource{}
		r.AddMethod("POST", resource.Method{})
		err := r.AddMethod("POST", resource.Method{})
		assertEqualError(t, err, resource.ErrorResourceMethodCollition)
	})
	t.Run("adding methods", func(t *testing.T) {
		r := resource.Resource{}
		err := r.AddMethod("POST", resource.Method{Summary: summary})
		assertNoErrorFatal(t, err)
		method, ok := r.GetMethod("post")
		if !ok {
			t.Errorf("Method was not found. got: %v", ok)
		}
		if method.Summary != summary {
			t.Errorf("got: %v want: %v", method.Summary, summary)
		}
	})
}

func TestGetMethod(t *testing.T) {
	t.Run("method found", func(t *testing.T) {
		r := resource.Resource{}
		err := r.AddMethod("POST", resource.Method{Summary: summary})
		assertNoErrorFatal(t, err)
		method, ok := r.GetMethod("post")
		if !ok {
			t.Errorf("Method was not found. got: %v", ok)
		}
		if method.Summary != summary {
			t.Errorf("got: %v want: %v", method.Summary, summary)
		}
	})
	t.Run("method not found", func(t *testing.T) {
		r := resource.Resource{}
		_, ok := r.GetMethod("post")
		if ok {
			t.Errorf("Not expecting method found: %v", ok)
		}
	})
}

func TestAddURIParamResource(t *testing.T) {
	r := resource.Resource{}
	paramResource, err := r.AddURIParamResource("{myParamID}", func(r *http.Request) string { return "" })
	assertNoErrorFatal(t, err)
	if !reflect.DeepEqual(paramResource, r.Resources[0]) {
		t.Errorf("got: %v want: %v", paramResource, r.Resources[0])
	}
}

func TestNewResource(t *testing.T) {
	t.Run("new resource", func(t *testing.T) {
		path := "/pet"
		r, err := resource.NewResource(path)
		assertNoErrorFatal(t, err)
		if r.Path != path {
			t.Errorf("got : %v want: %v", r.Path, path)
		}
	})
	t.Run("new resource with brackets", func(t *testing.T) {
		path := "/{pet}"
		_, err := resource.NewResource(path)
		assertEqualError(t, err, resource.ErrorResourceBracketsNotAllowed)
	})

}
