package chigenerator

import (
	"net/http"
	"strings"

	"github.com/ehsoc/rest"
	"github.com/go-chi/chi"
)

type ChiGenerator struct {
}

func (c ChiGenerator) GenerateServer(api rest.API) http.Handler {
	router := chi.NewMux()
	if strings.TrimSpace(api.BasePath) != "" {
		router.Route(api.BasePath, func(r chi.Router) {
			for _, resource := range api.Resources() {
				processResource(r, resource)
			}
		})
	}

	return router
}

func (c ChiGenerator) GetURIParam() func(*http.Request, string) string {
	return chi.URLParam
}

func processResource(r chi.Router, res rest.Resource) {
	r.Route("/"+res.Path(), func(r chi.Router) {
		for _, method := range res.Methods() {
			r.Method(method.HTTPMethod, "/", method)
		}
		for _, subRes := range res.Resources() {
			processResource(r, subRes)
		}
	})
}
