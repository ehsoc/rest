package docgen

import (
	"io"

	"github.com/ehsoc/resource"
)

//DocumentSpecGenerator will write an API Spec description in a determinate format
type APISpecGenerator interface {
	GenerateAPISpec(w io.Writer, restAPI resource.RestAPI)
}
