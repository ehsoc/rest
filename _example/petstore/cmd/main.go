package main

import (
	"net/http"
	"os"

	"github.com/ehsoc/resource"
	"github.com/ehsoc/resource/server_generators/chigenerator"
	"github.com/ehsoc/resource/test/petstore"
	"github.com/go-chi/chi"
	"github.com/swaggo/http-swagger"
)

func main() {
	f, _ := os.Create("oaiv2.json")
	api := petstore.GeneratePetStore()
	api.Host = "localhost:1323"
	api.GenerateSpec(f, &resource.OpenAPIV2SpecGenerator{})
	petServer := api.GenerateServer(chigenerator.ChiGenerator{})

	r := chi.NewRouter()

	r.Get("/doc/doc.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		api.GenerateSpec(w, &resource.OpenAPIV2SpecGenerator{})
	}))
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:1323/doc"), //The url pointing to API definition"

	))

	r.Mount("/", petServer)
	http.ListenAndServe(":1323", r)
}
