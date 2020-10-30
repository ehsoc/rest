package httputil_test

import (
	"reflect"
	"testing"

	"github.com/ehsoc/resource/httputil"
)

var negotiatorTests = []struct {
	accept    string
	renderers []httputil.MediaType
}{
	{"", []httputil.MediaType{{}}},
	{"application/json, application/xml", []httputil.MediaType{
		{Name: "application/json", Params: map[string]string{}},
		{Name: "application/xml", Params: map[string]string{}},
	}},
	{"application/json; indent=4, application/json, application/yaml, text/html, */*", []httputil.MediaType{
		{Name: "application/json", Params: map[string]string{"indent": "4"}},
		{Name: "application/json", Params: map[string]string{}},
		{Name: "application/yaml", Params: map[string]string{}},
		{Name: "text/html", Params: map[string]string{}},
		{Name: "*/*", Params: map[string]string{}},
	}},
}

func TestParseAccept(t *testing.T) {
	for _, tt := range negotiatorTests {
		t.Run("", func(t *testing.T) {
			got := httputil.ParseMediaTypes(tt.accept)
			if !reflect.DeepEqual(got, tt.renderers) {
				t.Errorf("got:%v want:%v", got, tt.renderers)
			}
		})
	}
}
