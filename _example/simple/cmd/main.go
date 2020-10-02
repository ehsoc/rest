package main

import (
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/ehsoc/resource"
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
	getCarOperation := func(i resource.Input, decoder encdec.Decoder) (interface{}, error) {
		carId, err := i.GetURIParam("carId")
		if err != nil {
			log.Println("error getting parameter: ", err)
			return Car{}, err
		}
		return Car{carId, "Mazda"}, nil
	}
	getCar := res.NewMethodOperation(
		res.OperationFunc(getCarOperation),
		res.NewResponse(200),
		res.NewResponse(404),
		true,
	)
	ct := res.NewHTTPContentTypeSelector()
	ct.Add("application/json", encdec.JSONEncoderDecoder{}, true)

	root := res.RestAPI{}
	root.BasePath = "/v1"
	root.Version = "v1"
	root.Title = "My simple car API"
	root.Resource("car", func(r *res.Resource) {
		carIdParam := res.NewURIParameter("carId", reflect.String)
		r.Resource("{carId}", func(r *res.Resource) {
			r.Get(getCar, ct).WithParameter(carIdParam)
		})
	})
	server := root.GenerateServer(chigenerator.ChiGenerator{})
	root.GenerateSpec(os.Stdout, &oaiv2.OpenAPIV2SpecGenerator{})
	http.ListenAndServe(":8080", server)
}
