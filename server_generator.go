package rest

import (
	"net/http"
)

// ServerGenerator interface describes the methods for generating the server and how to get the URI parameter.
type ServerGenerator interface {
	GenerateServer(API API) http.Handler
	GetURIParam() func(r *http.Request, key string) string
}
