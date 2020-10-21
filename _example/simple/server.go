package simple

import (
	"errors"
	"log"
	"net/http"
	"os"
	"reflect"

	res "github.com/ehsoc/resource"
	"github.com/ehsoc/resource/encdec"
	"github.com/ehsoc/resource/server_generator/chigenerator"
	"github.com/ehsoc/resource/specification_generator/oaiv2"
)

type Car struct {
	Id    string
	Brand string
}

func GenerateServer() http.Handler {
	getCarOperation := func(i res.Input) (body interface{}, success bool, err error) {
		carId, err := i.GetURIParam("carId")
		if err != nil {
			log.Println("error getting parameter: ", err)
			return Car{}, false, err
		}
		if carId == "error" {
			//Internal error trying to get the car. This will trigger a response code 500
			return nil, false, errors.New("Internal error")
		}
		if carId != "101" {
			//Car not found, success is false, no error. This will trigger a response code 404
			return nil, false, nil
		}
		//Car found, success true, error nil. This will trigger a response code 200
		return Car{carId, "Foo"}, true, nil
	}
	successResponse := res.NewResponse(200).WithOperationResultBody(Car{})
	failResponse := res.NewResponse(404)
	getCar := res.NewMethodOperation(res.OperationFunc(getCarOperation), successResponse).WithFailResponse(failResponse)

	rds := res.NewRenderers()
	rds.Add("application/json", encdec.JSONEncoderDecoder{}, true)

	root := res.RestAPI{}
	root.BasePath = "/v1"
	root.Version = "v1"
	root.Title = "My simple car API"
	root.Resource("car", func(r *res.Resource) {
		carIdParam := res.NewURIParameter("carId", reflect.String)
		r.Resource("{carId}", func(r *res.Resource) {
			r.Get(getCar, rds).WithParameter(carIdParam)
		})
	})
	server := root.GenerateServer(chigenerator.ChiGenerator{})
	root.GenerateSpec(os.Stdout, &oaiv2.OpenAPIV2SpecGenerator{})
	return server
}
