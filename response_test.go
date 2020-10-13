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

type MutableBodyStub struct {
	Code      int
	ErrorType string
	Message   string
}

var errorMessage = "this is my error"
var myError = "myError"

func (mr *MutableBodyStub) Mutate(v interface{}, success bool, err error) {
	mr.Code = 500
	mr.ErrorType = myError
	mr.Message = errorMessage
}

func TestMutateResponseBody(t *testing.T) {
	want := &MutableBodyStub{500, myError, errorMessage}
	mutableResponseBody := &MutableBodyStub{}
	resp := resource.NewResponse(500).WithMutableBody(mutableResponseBody)
	noChange := resp.Body()
	if !reflect.DeepEqual(noChange, mutableResponseBody) {
		t.Errorf("got: %#v \nwant: %#v", noChange, mutableResponseBody)
	}
	resp.Mutate(nil, false, nil)
	got := resp.Body()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %#v want: %#v", got, want)
	}
}

func TestWithOperationResultBody(t *testing.T) {
	resp := resource.NewResponse(200).WithOperationResultBody(Car{})
	want := Car{}
	got := resp.Body()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %#v want: %#v", got, want)
	}
	operationResult := Car{Brand: "Foo"}
	resp.Mutate(operationResult, false, nil)
	mutatedResponse := resp.Body()
	if !reflect.DeepEqual(mutatedResponse, operationResult) {
		t.Errorf("got: %#v want: %#v", got, want)
	}
}
