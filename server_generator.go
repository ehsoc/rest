package resource

import (
	"net/http"
)

type ServerGenerator interface {
	GenerateServer(restAPI RestAPI) http.Handler
}
