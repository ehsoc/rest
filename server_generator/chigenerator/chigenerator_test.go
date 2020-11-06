package chigenerator_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ehsoc/resource"
	"github.com/ehsoc/resource/encdec"
	"github.com/ehsoc/resource/server_generator/chigenerator"
	"github.com/ehsoc/resource/test/petstore"
)

type OperationStub struct {
	wasCall bool
	Pet     petstore.Pet
	PetID   string
}

func (o *OperationStub) Execute(i resource.Input) (interface{}, bool, error) {
	o.wasCall = true
	petID, _ := i.GetURIParam("petId")
	o.PetID = petID
	pet := petstore.Pet{}
	body, _ := i.GetBody()
	if body != nil && body != http.NoBody {
		i.BodyDecoder.Decode(body, &pet)
		o.Pet = pet
	}
	error, _ := i.GetQueryString("error")
	if error != "" {
		return nil, false, errors.New("Failed")
	}
	return o.Pet, true, nil
}

func TestGenerateServer(t *testing.T) {
	t.Run("get method", func(t *testing.T) {
		gen := chigenerator.ChiGenerator{}
		api := resource.RestAPI{}
		api.BasePath = "/v2"
		api.Host = "localhost"
		renderers := resource.NewRenderers()
		renderers.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		operation := &OperationStub{}
		getMethodOp := resource.NewMethodOperation(operation, resource.NewResponse(200)).WithFailResponse(resource.NewResponse(http.StatusNotFound))
		petResource := resource.NewResource("pet")
		petResource.Resource("{petId}", func(r *resource.Resource) {
			uriParam := resource.NewURIParameter("petId", reflect.String)
			r.Get(getMethodOp, renderers).WithParameter(uriParam)
		})
		myID := "101"
		api.AddResource(petResource)
		server := gen.GenerateServer(api)
		ctx := context.WithValue(context.Background(), resource.InputContextKey("uriparamfunc"), gen.GetURIParam())
		request, _ := http.NewRequest(http.MethodGet, "/v2/pet/"+myID, nil)
		request = request.WithContext(ctx)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Errorf("got: %v want: %v", response.Code, http.StatusOK)
		}
		if !operation.wasCall {
			t.Errorf("operation was not called")
		}
		if operation.PetID != myID {
			t.Errorf("got: %s want: %s", operation.PetID, myID)
		}
	})
	t.Run("post method", func(t *testing.T) {
		gen := chigenerator.ChiGenerator{}
		api := resource.RestAPI{}
		api.BasePath = "/v2"
		api.Host = "localhost"
		renderers := resource.NewRenderers()
		renderers.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		operation := &OperationStub{}
		postMethodOp := resource.NewMethodOperation(operation, resource.NewResponse(http.StatusCreated).WithBody(petstore.Pet{})).WithFailResponse(resource.NewResponse(http.StatusBadRequest))
		postMethod := resource.NewMethod(http.MethodPost, postMethodOp, renderers)
		petResource := resource.NewResource("pet")
		postMethod.RequestBody = resource.RequestBody{Description: "", Body: petstore.Pet{}}
		petResource.AddMethod(postMethod)

		api.AddResource(petResource)
		server := gen.GenerateServer(api)

		pet := petstore.Pet{Name: "Cat"}
		buf := new(bytes.Buffer)
		encoder := encdec.JSONEncoderDecoder{}
		encoder.Encode(buf, pet)

		request, _ := http.NewRequest(http.MethodPost, "/v2/pet", buf)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		if response.Code != http.StatusCreated {
			t.Errorf("got: %v want: %v", response.Code, http.StatusCreated)
		}
		if !operation.wasCall {
			t.Errorf("operation was not called")
		}
		if !reflect.DeepEqual(pet, operation.Pet) {
			t.Errorf("got: %v want: %v", pet, operation.Pet)
		}
	})
}

var testRoutes = []struct {
	route    string
	wantCode int
}{
	{"/v1/1", 404},
	{"/v1/1/2", 404},
	{"/v1/1/2/3", 200},
	{"/v1/1/2/3/4/5/1", 200},
}

func TestNestedRoutes(t *testing.T) {
	mo := resource.NewMethodOperation(&OperationStub{}, resource.NewResponse(http.StatusOK)).WithFailResponse(resource.NewResponse(500))
	ct := resource.NewRenderers()
	ct.Add("application/json", encdec.JSONEncoderDecoder{}, true)
	rootResource := resource.NewResource("1")
	rootResource.Resource("2", func(r *resource.Resource) {
		r.Resource("3", func(r *resource.Resource) {
			r.Get(mo, ct)
			r.Resource("4", func(r *resource.Resource) {
				r.Resource("5", func(r *resource.Resource) {
					r.Resource("1", func(r *resource.Resource) {
						r.Get(mo, ct)
					})
				})
			})
		})

	})
	api := resource.RestAPI{}
	api.BasePath = "/v1"
	api.AddResource(rootResource)
	server := api.GenerateServer(chigenerator.ChiGenerator{})
	for _, test := range testRoutes {
		t.Run(test.route, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, test.route, nil)
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)
			if response.Code != test.wantCode {
				t.Errorf("got: %v want: %v", response.Code, test.wantCode)
			}
		})

	}

}
