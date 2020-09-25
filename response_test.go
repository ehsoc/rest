package resource_test

import (
	"reflect"
	"testing"

	"github.com/ehsoc/resource"
)

func TestNewResponse(t *testing.T) {
	car := Car{ID: 1}
	description := "This description"
	t.Run("chain methods", func(t *testing.T) {
		r := resource.NewResponse(200).WithBody(car).WithDescription(description)
		if r.Body() == nil {
			t.Errorf("was not expecting a nil Body")
		}
		if !reflect.DeepEqual(r.Body(), car) {
			t.Errorf("got: %v want: %v", r.Body(), car)
		}
		if r.Description() != description {
			t.Errorf("got: %v want: %v", r.Description(), description)
		}
	})
	t.Run("chain methods should not work outside a chain", func(t *testing.T) {
		r := resource.NewResponse(200)
		if r.Body() != nil {
			t.Errorf("was expecting a nil Body")
		}
		r.WithBody(car)
		if r.Body() != nil {
			t.Errorf("was expecting a nil Body")
		}
		r.WithDescription(description)
		wantDesc := ""
		if r.Description() != wantDesc {
			t.Errorf("got: %v want: %v", r.Description(), wantDesc)
		}
	})

}
