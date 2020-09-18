package chigenerator_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ehsoc/resource"
	"github.com/ehsoc/resource/encdec"
	"github.com/ehsoc/resource/server_generators/chigenerator"
	"github.com/ehsoc/resource/test/petstore"
	"github.com/go-chi/chi"
)

type OperationStub struct {
	wasCall bool
	Pet     petstore.Pet
	PetId   string
}

func (o *OperationStub) Execute(r *http.Request, decoder encdec.Decoder) (interface{}, error) {
	o.wasCall = true
	o.PetId = chi.URLParam(r, "petId")
	pet := petstore.Pet{}
	if r.Body != nil && r.Body != http.NoBody {
		decoder.Decode(r.Body, &pet)
		o.Pet = pet
	}
	if r.URL.Query().Get("error") != "" {
		return nil, errors.New("Failed")
	}
	return o.Pet, nil
}

func TestGenerateServer(t *testing.T) {
	t.Run("get method", func(t *testing.T) {
		gen := chigenerator.ChiGenerator{}
		api := resource.RestAPI{}
		api.BasePath = "/v2"
		api.Host = "localhost"
		contentTypes := resource.NewHTTPContentTypeSelector(resource.Response{})
		contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		operation := &OperationStub{}
		getMethodOp := resource.NewMethodOperation(operation, resource.Response{http.StatusOK, nil, ""}, resource.Response{http.StatusNotFound, nil, ""}, true)
		getMethod := resource.NewMethod(http.MethodGet, getMethodOp, contentTypes)
		petResource, _ := resource.NewResourceWithURIParam("/pet/{petId}", "", reflect.String)
		getMethod.AddParameter(*petResource.GetURIParam())
		petResource.AddMethod(getMethod)
		myId := "101"
		api.Resources = append(api.Resources, &petResource)
		server := gen.GenerateServer(api)
		request, _ := http.NewRequest(http.MethodGet, "/v2/pet/"+myId, nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Errorf("got: %v want: %v", response.Code, http.StatusOK)
		}
		if !operation.wasCall {
			t.Errorf("operation was not called")
		}
		if operation.PetId != myId {
			t.Errorf("got: %s want: %s", operation.PetId, myId)
		}
	})
	t.Run("post method", func(t *testing.T) {
		gen := chigenerator.ChiGenerator{}
		api := resource.RestAPI{}
		api.BasePath = "/v2"
		api.Host = "localhost"
		contentTypes := resource.NewHTTPContentTypeSelector(resource.Response{http.StatusUnsupportedMediaType, nil, ""})
		contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		operation := &OperationStub{}
		postMethodOp := resource.NewMethodOperation(operation, resource.Response{http.StatusCreated, petstore.Pet{}, ""}, resource.Response{http.StatusBadRequest, nil, ""}, true)
		postMethod := resource.NewMethod(http.MethodPost, postMethodOp, contentTypes)
		petResource, _ := resource.NewResource("/pet")
		postMethod.RequestBody = resource.RequestBody{"", petstore.Pet{}}
		petResource.AddMethod(postMethod)

		api.Resources = append(api.Resources, &petResource)
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
	mo := resource.NewMethodOperation(&OperationStub{}, resource.Response{200, nil, ""}, resource.Response{500, nil, ""}, false)
	ct := resource.NewHTTPContentTypeSelector(resource.Response{415, nil, ""})
	ct.Add("application/json", encdec.JSONEncoderDecoder{}, true)
	method := resource.NewMethod(http.MethodGet, mo, ct)
	rootResource, _ := resource.NewResource("/1/2")
	r3, _ := resource.NewResource("/3")
	r5, _ := resource.NewResourceWithURIParam("/4/5/{petId}", "", reflect.String)
	r3.AddMethod(method)
	r5.AddMethod(method)
	r3.Resources = append(r3.Resources, &r5)
	rootResource.Resources = append(rootResource.Resources, &r3)
	api := resource.RestAPI{}
	api.BasePath = "/v1"
	api.Resources = append(api.Resources, &rootResource)
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
