package main

import (
	"net/http"
	"strings"

	"github.com/ehsoc/resource/server_generator/chigenerator"
	"github.com/ehsoc/resource/specification_generator/oaiv2"
	"github.com/ehsoc/resource/test/petstore"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/afero"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	api := petstore.GeneratePetStore()
	api.Host = "localhost:1323"
	api.Title = "My petstore"
	api.Version = "v1"
	petServer := api.GenerateServer(chigenerator.ChiGenerator{})

	r := chi.NewRouter()
	r.Use(middleware.DefaultLogger)

	r.Get("/doc.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		api.GenerateSpec(w, &oaiv2.OpenAPIV2SpecGenerator{})
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
