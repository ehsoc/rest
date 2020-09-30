package resource

import (
	"net/http"
)

type ServerGenerator interface {
	GenerateServer(restAPI RestAPI) http.Handler
	GetURIParam() func(r *http.Request, key string) string
}
