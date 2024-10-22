package httputil_test

import (
	"fmt"

	"github.com/ehsoc/rest/httputil"
)

func ExampleParseMediaTypes() {
	accept := "application/json; indent=4, application/xml"
	result := httputil.ParseMediaTypes(accept)

	for _, mt := range result {
		fmt.Println(mt.Name, mt.Params)
	}
	// Output:
	// application/json map[indent:4]
	// application/xml map[]
}
