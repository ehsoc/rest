package resource

import (
	"io"
)

//DocumentSpecGenerator will write an API Spec description in a determinate format
type APISpecGenerator interface {
	GenerateAPISpec(w io.Writer, restAPI RestAPI)
}
