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

	t.Run("parameter not defined", func(t *testing.T) {
		params := resource.Parameters{}
		_, err := params.GetParameter(resource.QueryParameter, "foo")
		if _, ok := err.(*resource.TypeErrorParameterNotDefined); !ok {
			t.Errorf("got: %v want: %T", err, resource.TypeErrorParameterNotDefined{})
		}
	})

}
