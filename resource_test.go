package resource_test

import (
	"fmt"
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
		r.AddMethod(resource.Method{})
	})
	t.Run("error on existing key", func(t *testing.T) {

		r := resource.Resource{}
		r.AddMethod(resource.Method{})
		r.AddMethod(resource.Method{})
	})
	t.Run("adding methods", func(t *testing.T) {
		r := resource.Resource{}
		err := r.AddMethod(resource.Method{HTTPMethod: http.MethodPost, Summary: summary})
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
		err := r.AddMethod(resource.Method{HTTPMethod: http.MethodPost, Summary: summary})
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

// func TestAddURIParamResource(t *testing.T) {
// 	r := resource.Resource{}
// 	paramResource, err := r.AddURIParamResource("{myParamID}", func(r *http.Request) string { return "" })
// 	assertNoErrorFatal(t, err)
// 	if !reflect.DeepEqual(paramResource, r.Resources[0]) {
// 		t.Errorf("got: %v want: %v", paramResource, r.Resources[0])
// 	}
// }

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

func TestResourceWithURIParam(t *testing.T) {
	t.Run("new resource with uri parameter", func(t *testing.T) {
		paramName := "petId"
		path := fmt.Sprintf("/pet/{%s}", paramName)
		r, err := resource.NewResourceWithURIParam(path, "", reflect.String)
		assertNoErrorFatal(t, err)
		if r.Path != path {
			t.Errorf("got : %v want: %v", r.Path, path)
		}
		URIParameter := r.GetURIParam()
		if URIParameter.Name != paramName {
			t.Errorf("got: %v want: %v ", URIParameter.Name, paramName)
		}
		if URIParameter.HTTPType != resource.URIParameter {
			t.Errorf("got: %v want: %v ", URIParameter.Type, resource.URIParameter)
		}
	})
	t.Run("new resource with uri parameter", func(t *testing.T) {
		paramName := "petId"
		path := fmt.Sprintf("/pet/{%s}", paramName)
		r, err := resource.NewResourceWithURIParam(path, "", reflect.String)
		assertNoErrorFatal(t, err)
		if r.Path != path {
			t.Errorf("got : %v want: %v", r.Path, path)
		}
	})
}

func TestResource(t *testing.T) {
	r, _ := resource.NewResource("/api")
	r.Resource("/pet", func(r *resource.Resource) {
		r.Resource("/1", nil)
	})

	if len(r.Resources) < 1 {
		t.Fatalf("not expecting empty resource slice")
	}
	if r.Resources[0].Path != "/pet" {
		t.Errorf("got: %v want: %v", r.Resources[0].Path, "/pet")
	}
	if len(r.Resources[0].Resources) < 1 {
		t.Fatalf("not expecting empty resource slice on %v", r.Resources[0].Path)
	}
	if r.Resources[0].Resources[0].Path != "/1" {
		t.Errorf("got: %v want: %v", r.Resources[0].Path, "/1")
	}
}
