package oaiv2_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ehsoc/rest"
	"github.com/ehsoc/rest/encdec"
	"github.com/ehsoc/rest/specification_generator/oaiv2"
	"github.com/ehsoc/rest/test/petstore"
	"github.com/go-openapi/spec"
	"github.com/nsf/jsondiff"
)

func TestNoEmptyResources(t *testing.T) {
	api := rest.API{}
	api.BasePath = "/v1"
	api.Version = "v1"
	api.Title = "My simple car API"
	api.Resource("car", func(r *rest.Resource) {
		carIDParam := rest.NewURIParameter("carId", reflect.String)
		r.Resource("{carId}", func(r *rest.Resource) {
			r.Get(rest.MethodOperation{}, rest.ContentTypes{}).WithParameter(carIDParam)
		})
	})

	gen := oaiv2.OpenAPIV2SpecGenerator{}
	generatedSpec := new(bytes.Buffer)
	decoder := json.NewDecoder(generatedSpec)
	gen.GenerateAPISpec(generatedSpec, api)
	gotSwagger := spec.Swagger{}
	decoder.Decode(&gotSwagger)
	if gotSwagger.Paths == nil {
		t.Fatal("not expecting nil Paths")
	}
	if len(gotSwagger.Paths.Paths) != 1 {
		t.Errorf("expecting just one resource, got: %v", len(gotSwagger.Paths.Paths))
	}
	wantPath := "/car/{carId}"
	_, ok := gotSwagger.Paths.Paths[wantPath]
	if !ok {
		t.Errorf("want: %v", wantPath)
	}
}

type TestStruct struct {
	CreatedAt time.Time `json:"created_at"`
}

func TestDateTime(t *testing.T) {
	api := rest.API{}
	api.Resource("test", func(r *rest.Resource) {
		mo := rest.NewMethodOperation(rest.OperationFunc(func(i rest.Input) (body interface{}, success bool, err error) {
			return nil, true, nil
		}), rest.NewResponse(200).WithBody(TestStruct{}))
		r.Post(mo, rest.NewContentTypes())
	})
	generatedSpec := new(bytes.Buffer)
	decoder := json.NewDecoder(generatedSpec)
	gen := oaiv2.OpenAPIV2SpecGenerator{}
	gen.GenerateAPISpec(generatedSpec, api)
	gotSwagger := spec.Swagger{}
	decoder.Decode(&gotSwagger)

	if !gotSwagger.Definitions["TestStruct"].Properties["created_at"].Type.Contains("string") {
		t.Errorf("Expecting type string")
	}

	if gotSwagger.Definitions["TestStruct"].Properties["created_at"].Format != "date-time" {
		t.Errorf("got: %v want: %v", gotSwagger.Definitions["TestStruct"].Properties["created_at"].Format, "date-time")
	}
}

func TestGenerateAPISpec(t *testing.T) {
	api := petstore.GeneratePetStore()
	gen := oaiv2.OpenAPIV2SpecGenerator{}
	generatedSpec := new(bytes.Buffer)
	decoder := json.NewDecoder(generatedSpec)
	gen.GenerateAPISpec(generatedSpec, api)
	gotSwagger := spec.Swagger{}
	decoder.Decode(&gotSwagger)
	wantSwagger := spec.Swagger{}
	petJSON := getPetJSON()
	err := json.Unmarshal(petJSON, &wantSwagger)
	assertNoErrorFatal(t, err)

	if gotSwagger.BasePath != wantSwagger.BasePath {
		t.Errorf("got: %v want: %v", gotSwagger.BasePath, wantSwagger.BasePath)
	}
	gotPetPath, ok := gotSwagger.Paths.Paths["/pet"]
	if !ok {
		t.Fatalf("Path not found")
	}
	wantPetPath, ok := wantSwagger.Paths.Paths["/pet"]
	if !ok {
		t.Fatalf("Path not found")
	}
	assertJSONStructEqual(t, gotPetPath.Post, wantPetPath.Post)
	assertJSONStructEqual(t, gotPetPath.Put, wantPetPath.Put)

	gotGetPetIDPath, ok := gotSwagger.Paths.Paths["/pet/{petId}"]
	if !ok {
		t.Fatalf("Path not found")
	}
	wantGetPetIDPath, ok := wantSwagger.Paths.Paths["/pet/{petId}"]
	if !ok {
		t.Fatalf("Path not found")
	}
	assertJSONStructEqual(t, gotGetPetIDPath.Get, wantGetPetIDPath.Get)
	// Delete
	assertJSONStructEqual(t, gotGetPetIDPath.Delete, wantGetPetIDPath.Delete)

	// Upload Image
	gotUploadImagePath, ok := gotSwagger.Paths.Paths["/pet/{petId}/uploadImage"]
	if !ok {
		t.Fatalf("Path not found")
	}
	wantUploadImagePath, ok := wantSwagger.Paths.Paths["/pet/{petId}/uploadImage"]
	if !ok {
		t.Fatalf("Path not found")
	}
	assertJSONStructEqual(t, gotUploadImagePath.Post, wantUploadImagePath.Post)

	// Find by status
	gotFindByStatusPath, ok := gotSwagger.Paths.Paths["/pet/findByStatus"]
	if !ok {
		t.Fatalf("Path not found")
	}
	wantFindByStatusPath, ok := wantSwagger.Paths.Paths["/pet/findByStatus"]
	if !ok {
		t.Fatalf("Path not found")
	}
	assertJSONStructEqual(t, gotFindByStatusPath.Get, wantFindByStatusPath.Get)

	t.Run("security definitions apiKey", func(t *testing.T) {
		if gotSwagger.SecurityDefinitions == nil {
			t.Fatal("SecurityDefinitions in nil")
		}
		gotAPIKeySchema, ok := gotSwagger.SecurityDefinitions["api_key"]
		if !ok {
			t.Fatal("expecting api_key in generated spec")
		}
		wantAPIKeySchema, ok := wantSwagger.SecurityDefinitions["api_key"]
		if !ok {
			t.Fatal("expecting api_key in test fixture")
		}
		assertJSONStructEqual(t, gotAPIKeySchema, wantAPIKeySchema)
	})
	t.Run("security definitions oauth2", func(t *testing.T) {
		if gotSwagger.SecurityDefinitions == nil {
			t.Fatal("SecurityDefinitions in nil")
		}
		gotOAuth2Schema, ok := gotSwagger.SecurityDefinitions["petstore_auth"]
		if !ok {
			t.Fatal("expecting petstore_auth in generated spec")
		}
		wantOAuth2Schema, ok := wantSwagger.SecurityDefinitions["petstore_auth"]
		if !ok {
			t.Fatal("expecting petstore_auth in test fixture")
		}
		assertJSONStructEqual(t, gotOAuth2Schema, wantOAuth2Schema)
	})
	t.Run("security definitions basic", func(t *testing.T) {
		if gotSwagger.SecurityDefinitions == nil {
			t.Fatal("SecurityDefinitions in nil")
		}
		gotSecuritySchema, ok := gotSwagger.SecurityDefinitions["basicSecurity"]
		if !ok {
			t.Fatal("expecting basicSecurity in generated spec")
		}
		wantSecuritySchema, ok := wantSwagger.SecurityDefinitions["basicSecurity"]
		if !ok {
			t.Fatal("expecting basicSecurity in test fixture")
		}
		assertJSONStructEqual(t, gotSecuritySchema, wantSecuritySchema)
	})
}

func assertJSONStructEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	gotJSON, err := json.MarshalIndent(got, " ", "  ")
	if err != nil {
		t.Fatalf("Not expecting error: %v", err)
	}
	wantJSON, err := json.MarshalIndent(want, " ", "  ")
	if err != nil {
		t.Fatalf("Not expecting error: %v", err)
	}

	if !reflect.DeepEqual(gotJSON, wantJSON) {
		opts := jsondiff.DefaultConsoleOptions()
		opts.PrintTypes = false
		_, result := jsondiff.Compare(gotJSON, wantJSON, &opts)
		t.Errorf("Expecting equal, diff: %s", result)
	}
}

func getPetJSON() []byte {
	jsonFile, err := os.Open("../../test/fixtures/petstore_oav2.json")
	if err != nil {
		log.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}

func assertNoErrorFatal(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Was not expecting error: %v", err)
	}
}

var SecOpStub = rest.SecurityOperation{
	Authenticator: rest.AuthenticatorFunc(func(i rest.Input) rest.AuthError {
		return nil
	}),
	FailedAuthenticationResponse: rest.NewResponse(401),
	FailedAuthorizationResponse:  rest.NewResponse(403),
}

func TestSecurityTwoSchemes(t *testing.T) {
	api := rest.API{}
	api.Resource("one", func(r *rest.Resource) {
		mo := rest.NewMethodOperation(rest.OperationFunc(func(i rest.Input) (interface{}, bool, error) {
			return nil, true, nil
		}), rest.NewResponse(200))
		ct := rest.NewContentTypes()
		ct.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		apiKeySchema := rest.NewSecurityScheme("api-key", rest.APIKeySecurityType, SecOpStub)
		apiKeySchema.AddParameter(rest.NewHeaderParameter("X-API-KEY", reflect.String))
		IDKeySchema := rest.NewSecurityScheme("id-key", rest.APIKeySecurityType, SecOpStub)
		IDKeySchema.AddParameter(rest.NewHeaderParameter("X-ID-KEY", reflect.String))

		r.Get(mo, ct).
			WithSecurity(
				rest.AddScheme(apiKeySchema),
				rest.AddScheme(IDKeySchema),
				rest.Enforce(),
			)
	})
	gen := oaiv2.OpenAPIV2SpecGenerator{}
	generatedSpec := new(bytes.Buffer)
	decoder := json.NewDecoder(generatedSpec)
	gen.GenerateAPISpec(generatedSpec, api)
	gotSwagger := spec.Swagger{}
	decoder.Decode(&gotSwagger)
	sec := gotSwagger.Paths.Paths["/one"].Get.Security
	if len(sec) != 1 {
		t.Fatalf("got: %d, expecting 1 in security slice", len(sec))
	}
	if _, ok := sec[0]["api-key"]; !ok {
		t.Errorf("expecting api-key map key")
	}
	if _, ok := sec[0]["id-key"]; !ok {
		t.Errorf("expecting id-key map key")
	}
}
