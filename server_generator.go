package resource

import (
	"net/http"
)

type GetURIParamFunc func(r *http.Request, key string) string

type ServerGenerator interface {
	GenerateServer(restAPI RestAPI) http.Handler
	GetURIParam() GetURIParamFunc
}
