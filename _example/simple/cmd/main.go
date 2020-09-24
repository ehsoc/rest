package main

import (
	"net/http"
	"os"
	"reflect"

	res "github.com/ehsoc/resource"
	"github.com/ehsoc/resource/document_generator/oaiv2"
	"github.com/ehsoc/resource/encdec"
	"github.com/ehsoc/resource/server_generator/chigenerator"
)

type Car struct {
	Id    string
	Brand string
}

func main() {
	getCarOperation := func(r *http.Request, decoder encdec.Decoder) (interface{}, error) {
		return Car{"1", "Mazda"}, nil
	}
	getCar := res.NewMethodOperation(
		res.OperationFunc(getCarOperation),
		res.Response{200, nil, ""},
		res.Response{404, nil, ""},
		true,
	)
	ct := res.NewHTTPContentTypeSelector(res.DefaultUnsupportedMediaResponse)
	ct.Add("application/json", encdec.JSONEncoderDecoder{}, true)

	root := res.RestAPI{}
	root.BasePath = "/v2"
	root.Resource("car", func(r *res.Resource) {
		carIdParam := res.NewURIParameter("carId", reflect.String)
		r.Resource("{carId}", func(r *res.Resource) {
			r.Get(getCar, ct).WithParameter(*carIdParam)
		})
	})
	server := root.GenerateServer(chigenerator.ChiGenerator{})
	root.GenerateSpec(os.Stdout, &oaiv2.OpenAPIV2SpecGenerator{})
	http.ListenAndServe(":8080", server)

}
