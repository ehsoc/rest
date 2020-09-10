package chigenerator_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	entity  interface{}
	Pet     petstore.Pet
	PetId   string
}

func (o *OperationStub) Execute(body io.ReadCloser, params url.Values, decoder encdec.Decoder) (interface{}, error) {
	o.wasCall = true
	o.PetId = params.Get("petId")
	pet := petstore.Pet{}
	if body != nil && body != http.NoBody {
		decoder.Decode(body, &pet)
		o.entity = pet
	}
	if params.Get("error") != "" {
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
		petResource, _ := resource.NewResourceWithURIParam("/pet/{petId}", resource.GetterFunc(func(r *http.Request) string {
			return chi.URLParam(r, "petId")
		}), "", reflect.String)
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

}
