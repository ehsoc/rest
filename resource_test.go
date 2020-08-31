package resource_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ehsoc/resource"
	"github.com/ehsoc/resource/encdec"
)

type Pet struct {
	ID        int64    `json:"id,omitempty" xml:"id,omitempty"`
	Name      string   `json:"name" xml:"name"`
	PhotoUrls []string `json:"photoUrls" xml:"photoUrl"`
	Status    string   `json:"status,omitempty" xml:"status,omitempty"`
}

func TestOpenAPIDocument(t *testing.T) {
	jsonString := `{
		"/pet": {
			"post": {
				"tags": [
					"pet"
				],
				"summary": "Add a new pet to the store",
				"description": "",
				"operationId": "addPet",
				"consumes": [
					"application/json",
					"application/xml"
				],
				"produces": [
					"application/xml",
					"application/json"
				],
				"parameters": [{
					"in": "body",
					"name": "body",
					"description": "Pet object that needs to be added to the store",
					"required": true,
					"schema": {
						"$ref": "#/definitions/Pet"
					}
				}],
				"responses": {
					"405": {
						"description": "Invalid input"
					}
				}
			}
		}
	}`
	successResponse := resource.Response{http.StatusCreated, nil}
	failedResponse := resource.Response{http.StatusBadRequest, nil}
	mo := resource.NewMethodOperation(Pet{}, nil, successResponse, failedResponse, nil, true, true)
	cts := resource.NewHTTPContentTypeSelector(resource.Response{http.StatusUnsupportedMediaType, nil})
	cts.Add("application/json", encdec.JSONEncoderDecoder{}, true)
	method := resource.NewMethod(mo, cts)

	res := resource.NewResource("/pets")
	res.AddMethod(http.MethodPost, method)
	gotString := res.GenerateOpenApiJSON()
	want := mustCompactJSON(t, jsonString)
	got := mustCompactJSON(t, gotString)
	if got != want {
		t.Errorf("\ngot: %s \nwant: %s", got, want)
	}

}

func mustCompactJSON(t *testing.T, jsonString string) string {
	bufWant := new(bytes.Buffer)
	err := json.Compact(bufWant, []byte(jsonString))
	if err != nil {
		t.Fatalf("Error compacting json: %v \nerror: %v", jsonString, err)
	}
	return bufWant.String()
}
