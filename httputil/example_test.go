package httputil_test

import (
	"fmt"

	"github.com/ehsoc/restapigen/httputil"
)

func ExampleParseContentType() {
	accept := "application/json; indent=4, application/xml"
	result := httputil.ParseContentType(accept)
	for _, mt := range result {
		fmt.Println(mt.Name, mt.Params)
	}
	// Output:
	// application/json map[indent:4]
	// application/xml map[]
}
