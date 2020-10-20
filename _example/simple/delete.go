package simple

import (
	"github.com/ehsoc/resource"
)

func t() {

	api := resource.RestAPI{}
	api.BasePath = "/v2"
	api.Host = "localhost"

	api.Resource("car", func(r *resource.Resource) {
		r.Resource("findMatch", func(r *resource.Resource) {
		})
		r.Resource("{carId}", func(r *resource.Resource) {
		})
	})

	api.Resource("ping", func(r *resource.Resource) {
	})

	api.Resource("user", func(r *resource.Resource) {
		r.Resource("signOut", func(r *resource.Resource) {
		})
	})

}
