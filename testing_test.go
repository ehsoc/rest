package resource_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ehsoc/resource"
	"github.com/ehsoc/resource/encdec"
)

var octetMimeType = "application/octet-stream"
var jsonMimeType = "application/json"

type Pet struct {
	// id
	ID int64 `json:"id,omitempty"`
	// name
	// Required: true
	Name string `json:"name"`
	// photo urls
	// Required: true
	PhotoUrls []string `json:"photoUrls" xml:"photoUrl"`
	// pet status in the store
	// Enum: [available pending sold]
	Status string `json:"status,omitempty"`
	// tags
	Tags []Tag `json:"tags" xml:"tag"`
}

// Tag tag
//
// swagger:model Tag
type Tag struct {
	// id
	ID int64 `json:"id,omitempty"`
	// name
	Name string `json:"name,omitempty"`
}

func assertResponseCode(t *testing.T, r *httptest.ResponseRecorder, want int) {
	t.Helper()
	if r.Code != want {
		t.Errorf("got: %v want:%v", r.Code, want)
	}
}

func assertTrue(t *testing.T, got bool) {
	t.Helper()
	if !got {
		t.Errorf("expecting to be true, got : %v", got)
	}
}

func assertFalse(t *testing.T, got bool) {
	t.Helper()
	if got {
		t.Errorf("expecting to be false, got : %v", got)
	}
}

func assertResponseContentType(t *testing.T, response *httptest.ResponseRecorder, mimeType string) {
	t.Helper()
	if response.Header().Get("Content-type") != mimeType {
		t.Errorf("got:%s want:%s", response.Header().Get("Content-type"), mimeType)
	}
}

func assertNoErrorFatal(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Was not expecting error: %v", err)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("Was expecting error")
	}
}

func assertEqualError(t *testing.T, err, want error) {
	if err != want {
		t.Errorf("got:%v want:%v", err, want)
	}
}

func generatePetStore() resource.RestAPI {
	getIdFunc := func(r *http.Request) string {
		return "id"
	}
	api := resource.RestAPI{}
	api.BasePath = "/v2"
	api.Host = "localhost"
	pets := resource.NewResource("/pet")
	contentTypes := resource.NewHTTPContentTypeSelector(resource.Response{})
	contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
	contentTypes.Add("application/xml", encdec.JSONEncoderDecoder{}, false)
	createMethodOperation := resource.NewMethodOperation(Pet{}, nil, resource.Response{201, nil}, resource.Response{405, nil}, getIdFunc, true, true)
	createPetMethod := resource.NewMethod(createMethodOperation, contentTypes)
	createPetMethod.Summary = "Add a new pet to the store"
	createPetMethod.Request.Description = "Pet object that needs to be added to the store"
	pets.AddMethod(http.MethodPost, createPetMethod)
	api.Resources = append(api.Resources, pets)
	return api
}
