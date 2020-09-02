package resource

import (
	"io"
)

//RestAPI is the root of a REST API abstraction.
//It is responsable document generation like output Open API v2 json generation and
//Server handler generation
type RestAPI struct {
	ID        string
	Host      string
	BasePath  string
	Resources []Resource
}

func (r RestAPI) GenerateAPISpec(w io.Writer, docGenerator APISpecGenerator) {
	docGenerator.GenerateAPISpec(w, r)
}