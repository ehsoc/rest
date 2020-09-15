package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/ehsoc/resource"
	"github.com/ehsoc/resource/server_generators/chigenerator"
	"github.com/ehsoc/resource/test/petstore"
	"github.com/go-chi/chi"
	"github.com/spf13/afero"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	f, _ := os.Create("oaiv2.json")
	api := petstore.GeneratePetStore()
	api.Host = "localhost:1323"
	api.GenerateSpec(f, &resource.OpenAPIV2SpecGenerator{})
	petServer := api.GenerateServer(chigenerator.ChiGenerator{})

	r := chi.NewRouter()

	r.Get("/doc.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		api.GenerateSpec(w, &resource.OpenAPIV2SpecGenerator{})
	}))
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:1323/doc.json"), //The url pointing to API definition"

	))
	httpFs := afero.NewHttpFs(petstore.PetStore.InMemoryFs)
	FileServer(r, "/files", httpFs.Dir("files"))

	r.Mount("/", petServer)
	http.ListenAndServe(":1323", r)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
