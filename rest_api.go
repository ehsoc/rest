package resource

import (
	"io"

	"github.com/ehsoc/resource/docgen"
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

func GenerateAPISpec(w io.Writer, docGenerator docgen.APISpecGenerator) {
	return docGenerator.GenerateAPISpec(w, docGenerator)
}
