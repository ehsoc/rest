package rest_test

import (
	"reflect"
	"testing"

	"github.com/ehsoc/rest"
	"github.com/ehsoc/rest/encdec"
)

var testParameters = []struct {
	parameter rest.Parameter
}{
	{rest.NewQueryParameter("foo", reflect.String)},
	{rest.NewURIParameter("foo", reflect.String)},
	{rest.NewFileParameter("foo")},
	{rest.NewFormDataParameter("foo", reflect.String, encdec.JSONEncoderDecoder{})},
	{rest.NewHeaderParameter("foo", reflect.String)},
}

func TestGetParameter(t *testing.T) {
	for _, tt := range testParameters {
		params := rest.ParameterCollection{}
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
		params := rest.ParameterCollection{}
		_, err := params.GetParameter(rest.QueryParameter, "foo")
		if _, ok := err.(*rest.TypeErrorParameterNotDefined); !ok {
			t.Errorf("got: %v want: %T", err, rest.TypeErrorParameterNotDefined{})
		}
	})
	t.Run("parameter not defined (empty parameter type)", func(t *testing.T) {
		params := rest.ParameterCollection{}
		params.AddParameter(rest.NewURIParameter("", reflect.Int))
		_, err := params.GetParameter(rest.QueryParameter, "foo")
		if _, ok := err.(*rest.TypeErrorParameterNotDefined); !ok {
			t.Errorf("got: %v want: %T", err, rest.TypeErrorParameterNotDefined{})
		}
	})
}

func TestWithBody(t *testing.T) {
	t.Run("set body", func(t *testing.T) {
		car := Car{}
		p := rest.NewFormDataParameter("car", reflect.Struct, encdec.JSONDecoder{}).WithBody(car)
		if !reflect.DeepEqual(p.Body, car) {
			t.Errorf("got: %v want: %v", p.Body, car)
		}
	})
	t.Run("method has not a pointer receiver", func(t *testing.T) {
		car := Car{}
		p := rest.NewFormDataParameter("car", reflect.Struct, encdec.JSONDecoder{})
		p.WithBody(car)
		// WithBody shouldn't change p properties
		if p.Body != nil {
			t.Errorf("got: %v want: %v", p.Body, nil)
		}
	})
}

func TestAsOptional(t *testing.T) {
	t.Run("as optional", func(t *testing.T) {
		p := rest.NewURIParameter("foo", reflect.Int).AsOptional()
		if p.Required {
			t.Errorf("expecting to be false")
		}
	})
	t.Run("method has not a pointer receiver", func(t *testing.T) {
		p := rest.NewURIParameter("foo", reflect.Int)
		p.AsOptional()
		if !p.Required {
			t.Errorf("expecting to be true")
		}
	})
}

func TestAsRequired(t *testing.T) {
	t.Run("as required", func(t *testing.T) {
		p := rest.NewQueryParameter("foo", reflect.String).AsRequired()
		if !p.Required {
			t.Errorf("expecting to be true")
		}
	})
	t.Run("method has not a pointer receiver", func(t *testing.T) {
		p := rest.NewQueryParameter("foo", reflect.String)
		p.AsRequired()
		if p.Required {
			t.Errorf("expecting to be false")
		}
	})
}

func TestWithDescription(t *testing.T) {
	description := "My description"
	t.Run("set description", func(t *testing.T) {
		p := rest.NewQueryParameter("foo", reflect.String).WithDescription(description)
		assertStringEqual(t, p.Description, description)
	})
	t.Run("method has not a pointer receiver", func(t *testing.T) {
		p := rest.NewQueryParameter("foo", reflect.String)
		p.WithDescription(description)
		assertStringEqual(t, p.Description, "")
	})
}
