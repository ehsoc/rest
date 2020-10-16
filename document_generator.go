package resource

import (
	"io"
)

// APISpecGenerator is the interface implemented by types that transform a RestAPI type into an REST API specification,
// in a specific format.
type APISpecGenerator interface {
	GenerateAPISpec(w io.Writer, restAPI RestAPI)
}
