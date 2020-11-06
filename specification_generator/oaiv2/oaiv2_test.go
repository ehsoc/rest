package oaiv2_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/ehsoc/resource"
	"github.com/ehsoc/resource/specification_generator/oaiv2"
	"github.com/ehsoc/resource/test/petstore"
	"github.com/go-openapi/spec"
	"github.com/nsf/jsondiff"
)

func TestNoEmptyResources(t *testing.T) {
	api := resource.RestAPI{}
	api.BasePath = "/v1"
	api.Version = "v1"
	api.Title = "My simple car API"
	api.Resource("car", func(r *resource.Resource) {
		carIDParam := resource.NewURIParameter("carId", reflect.String)
		r.Resource("{carId}", func(r *resource.Resource) {
			r.Get(resource.MethodOperation{}, resource.Renderers{}).WithParameter(carIDParam)
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
	//Delete
	assertJSONStructEqual(t, gotGetPetIDPath.Delete, wantGetPetIDPath.Delete)

	//Upload Image
	gotUploadImagePath, ok := gotSwagger.Paths.Paths["/pet/{petId}/uploadImage"]
	if !ok {
		t.Fatalf("Path not found")
	}
	wantUploadImagePath, ok := wantSwagger.Paths.Paths["/pet/{petId}/uploadImage"]
	if !ok {
		t.Fatalf("Path not found")
	}
	assertJSONStructEqual(t, gotUploadImagePath.Post, wantUploadImagePath.Post)

	//Find by status
	gotFindByStatusPath, ok := gotSwagger.Paths.Paths["/pet/findByStatus"]
	if !ok {
		t.Fatalf("Path not found")
	}
	wantFindByStatusPath, ok := wantSwagger.Paths.Paths["/pet/findByStatus"]
	if !ok {
		t.Fatalf("Path not found")
	}
	//fmt.Printf("%#v", gotFindByStatusPath)
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
