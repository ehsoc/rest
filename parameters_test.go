package resource_test

import (
	"reflect"
	"testing"

	"github.com/ehsoc/resource"
	"github.com/ehsoc/resource/encdec"
)

var testParameters = []struct {
	parameter resource.Parameter
}{
	{resource.NewQueryParameter("foo")},
	{resource.NewURIParameter("foo", reflect.String)},
	{resource.NewFileParameter("foo")},
	{resource.NewFormDataParameter("foo", reflect.String, encdec.JSONEncoderDecoder{})},
	{resource.NewHeaderParameter("foo", reflect.String)},
}

func TestGetParameter(t *testing.T) {
	for _, tt := range testParameters {
		params := resource.Parameters{}
		params.AddParameter(tt.parameter)
		t.Run(string(tt.parameter.HTTPType), func(t *testing.T) {
			gotParam, err := params.GetParameter(tt.parameter.HTTPType, tt.parameter.Name)
			assertNoErrorFatal(t, err)
			if !reflect.DeepEqual(gotParam, tt.parameter) {
				t.Errorf("got: %v want: %v", gotParam, tt.parameter)
			}
		})

	}
	t.Run("parameter nil collection", func(t *testing.T) {
		params := resource.Parameters{}
		_, err := params.GetParameter(resource.QueryParameter, "foo")
		if _, ok := err.(*resource.TypeErrorParameterNotDefined); !ok {
			t.Errorf("got: %v want: %T", err, resource.TypeErrorParameterNotDefined{})
		}
	})
	t.Run("parameter not defined (empty parameter type)", func(t *testing.T) {
		params := resource.Parameters{}
		params.AddParameter(resource.NewURIParameter("", reflect.Int))
		_, err := params.GetParameter(resource.QueryParameter, "foo")
		if _, ok := err.(*resource.TypeErrorParameterNotDefined); !ok {
			t.Errorf("got: %v want: %T", err, resource.TypeErrorParameterNotDefined{})
		}
	})
}

func TestWithBody(t *testing.T) {
	t.Run("set body", func(t *testing.T) {
		car := Car{}
		p := resource.NewFormDataParameter("car", reflect.Struct, encdec.JSONDecoder{}).WithBody(car)
		if !reflect.DeepEqual(p.Body, car) {
			t.Errorf("got: %v want: %v", p.Body, car)
		}
	})
	t.Run("method has not a pointer receiver", func(t *testing.T) {
		car := Car{}
		p := resource.NewFormDataParameter("car", reflect.Struct, encdec.JSONDecoder{})
		p.WithBody(car)
		//WithBody shouldn't change p properties
		if p.Body != nil {
			t.Errorf("got: %v want: %v", p.Body, nil)
		}
	})
}

func TestAsOptional(t *testing.T) {
	t.Run("as optional", func(t *testing.T) {
		p := resource.NewURIParameter("foo", reflect.Int).AsOptional()
		if p.Required {
			t.Errorf("expecting to be false")
		}
	})
	t.Run("method has not a pointer receiver", func(t *testing.T) {
		p := resource.NewURIParameter("foo", reflect.Int)
		p.AsOptional()
		if !p.Required {
			t.Errorf("expecting to be true")
		}
	})
}

func TestAsRequired(t *testing.T) {
	t.Run("as required", func(t *testing.T) {
		p := resource.NewQueryParameter("foo").AsRequired()
		if !p.Required {
			t.Errorf("expecting to be true")
		}
	})
	t.Run("method has not a pointer receiver", func(t *testing.T) {
		p := resource.NewQueryParameter("foo")
		p.AsRequired()
		if p.Required {
			t.Errorf("expecting to be false")
		}
	})
}

func TestWithDescription(t *testing.T) {
	description := "My description"
	t.Run("set description", func(t *testing.T) {
		p := resource.NewQueryParameter("foo").WithDescription(description)
		assertStringEqual(t, p.Description, description)
	})
	t.Run("method has not a pointer receiver", func(t *testing.T) {
		p := resource.NewQueryParameter("foo")
		p.WithDescription(description)
		assertStringEqual(t, p.Description, "")
	})
}
