package chigenerator

import (
	"log"
	"net/http"
	"strings"

	"github.com/ehsoc/resource"
	"github.com/go-chi/chi"
	"github.com/go-chi/docgen"
)

type ChiGenerator struct {
}

func (c ChiGenerator) GenerateServer(api resource.RestAPI) http.Handler {
	router := chi.NewMux()
	if strings.TrimSpace(api.BasePath) != "" {
		router.Route(api.BasePath, func(r chi.Router) {
			for _, resource := range api.Resources {
				processResource(r, resource)
			}
		})
	}
	log.Println("Generated routes:")
	docgen.PrintRoutes(router)
	return router
}

func processResource(r chi.Router, res *resource.Resource) {
	r.Route(res.Path, func(r chi.Router) {
		for httpMethod, method := range res.Methods {
			r.Method(httpMethod, "/", method)
		}
		for _, subRes := range res.Resources {
			processResource(r, subRes)
		}
	})
}
