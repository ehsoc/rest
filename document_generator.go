package resource

import (
	"io"
)

// APISpecGenerator is the interface implemented by types that transforms a RestAPI into a API specification,
// in a specific format.
type APISpecGenerator interface {
	GenerateAPISpec(w io.Writer, restAPI RestAPI)
}
