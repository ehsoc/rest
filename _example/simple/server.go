package simple

import (
	"errors"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/ehsoc/rest"
	"github.com/ehsoc/rest/encdec"
	"github.com/ehsoc/rest/server_generator/chigenerator"
	"github.com/ehsoc/rest/specification_generator/oaiv2"
)

type Car struct {
	ID    string
	Brand string
}

func GenerateServer() http.Handler {
	getCarOperation := func(i rest.Input) (body interface{}, success bool, err error) {
		carID, err := i.GetURIParam("carID")

		if err != nil {
			log.Println("error getting parameter: ", err)
			return Car{}, false, err
		}

		if carID == "error" {
			// Internal error trying to get the car. This will trigger a response code 500
			return nil, false, errors.New("Internal error")
		}

		if carID != "101" {
			// Car not found, success is false, no error. This will trigger a response code 404
			return nil, false, nil
		}

		// Car found, success true, error nil. This will trigger a response code 200
		return Car{carID, "Foo"}, true, nil
	}
	// Responses
	successResponse := rest.NewResponse(200).WithOperationResultBody(Car{})
	failResponse := rest.NewResponse(404)
	// Method Operation
	getCar := rest.NewMethodOperation(rest.OperationFunc(getCarOperation), successResponse).WithFailResponse(failResponse)
	// Renderers
	rds := rest.NewRenderers()
	rds.Add("application/json", encdec.JSONEncoderDecoder{}, true)

	api := rest.API{}
	api.BasePath = "/v1"
	api.Version = "v1"
	api.Title = "My simple car API"
	api.Resource("car", func(r *rest.Resource) {
		carIDParam := rest.NewURIParameter("carID", reflect.String)
		r.Resource("{carID}", func(r *rest.Resource) {
			r.Get(getCar, rds).WithParameter(carIDParam)
		})
	})
	// Generating OpenAPI v2 specification to standard output
	api.GenerateSpec(os.Stdout, &oaiv2.OpenAPIV2SpecGenerator{})

	// Generating server routes
	server := api.GenerateServer(chigenerator.ChiGenerator{})

	return server
}
