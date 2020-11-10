package rest

import (
	"io"
)

// APISpecGenerator is the interface implemented by types that transform a API type into an REST API specification,
// in a specific format and writing it to w.
type APISpecGenerator interface {
	GenerateAPISpec(w io.Writer, API API)
}
