package httputil_test

import (
	"reflect"
	"testing"

	"github.com/ehsoc/restapigen/httputil"
)

var negotiatorTests = []struct {
	accept       string
	contentTypes []httputil.MediaType
}{
	{"", []httputil.MediaType{httputil.MediaType{}}},
	{"application/json, application/xml", []httputil.MediaType{
		{"application/json", map[string]string{}},
		{"application/xml", map[string]string{}},
	}},
	{"application/json; indent=4, application/json, application/yaml, text/html, */*", []httputil.MediaType{
		{"application/json", map[string]string{"indent": "4"}},
		{"application/json", map[string]string{}},
		{"application/yaml", map[string]string{}},
		{"text/html", map[string]string{}},
		{"*/*", map[string]string{}},
	}},
}

func TestParseAccept(t *testing.T) {
	for _, tt := range negotiatorTests {
		t.Run("", func(t *testing.T) {
			got := httputil.ParseContentType(tt.accept)
			if !reflect.DeepEqual(got, tt.contentTypes) {
				t.Errorf("got:%v want:%v", got, tt.contentTypes)
			}
		})
	}
}
