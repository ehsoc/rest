package oaiv2_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/ehsoc/resource/document_generators/oaiv2"
	"github.com/ehsoc/resource/test/petstore"
	"github.com/go-openapi/spec"
	"github.com/nsf/jsondiff"
)

func TestGenerateAPISpec(t *testing.T) {
	api := petstore.GeneratePetStore()
	gen := oaiv2.OpenAPIV2SpecGenerator{}
	generatedSpec := new(bytes.Buffer)
	decoder := json.NewDecoder(generatedSpec)
	gen.GenerateAPISpec(generatedSpec, api)
	gotSwagger := spec.Swagger{}
	decoder.Decode(&gotSwagger)
	wantSwagger := spec.Swagger{}
	petJson := getPetJson()
	err := json.Unmarshal(petJson, &wantSwagger)
	assertNoErrorFatal(t, err)

	//TODO check definitions

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
	assertOAv2OperationEqual(t, gotPetPath.Post, wantPetPath.Post)

	gotGetPetIDPath, ok := gotSwagger.Paths.Paths["/pet/{petId}"]
	if !ok {
		t.Fatalf("Path not found")
	}
	wantGetPetIDPath, ok := wantSwagger.Paths.Paths["/pet/{petId}"]
	if !ok {
		t.Fatalf("Path not found")
	}
	assertOAv2OperationEqual(t, gotGetPetIDPath.Get, wantGetPetIDPath.Get)
	//Delete
	assertOAv2OperationEqual(t, gotGetPetIDPath.Delete, wantGetPetIDPath.Delete)

	//Upload Image
	gotUploadImagePath, ok := gotSwagger.Paths.Paths["/pet/{petId}/uploadImage"]
	if !ok {
		t.Fatalf("Path not found")
	}
	wantUploadImagePath, ok := wantSwagger.Paths.Paths["/pet/{petId}/uploadImage"]
	if !ok {
		t.Fatalf("Path not found")
	}
	assertOAv2OperationEqual(t, gotUploadImagePath.Post, wantUploadImagePath.Post)

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
	assertOAv2OperationEqual(t, gotFindByStatusPath.Get, wantFindByStatusPath.Get)

}

func assertJsonSchemaEqual(t *testing.T, got, want string) {
	gotJson := spec.Schema{}
	err := json.Unmarshal([]byte(got), &gotJson)
	if err != nil {
		t.Fatalf("Not expecting error: %v", err)
	}
	wantJson := spec.Schema{}
	err = json.Unmarshal([]byte(want), &wantJson)
	if err != nil {
		t.Fatalf("Not expecting error: %v", err)
	}
	if !reflect.DeepEqual(gotJson, wantJson) {
		t.Errorf("\ngot: %v \nwant: %v", got, want)
	}
}

func assertOAv2OperationEqual(t *testing.T, got, want *spec.Operation) {
	t.Helper()
	gotJson, err := json.MarshalIndent(got, " ", "  ")
	wantJson, err := json.MarshalIndent(want, " ", "  ")
	if err != nil {
		t.Fatalf("Not expecting error: %v", err)
	}
	opts := jsondiff.DefaultConsoleOptions()
	opts.PrintTypes = false
	_, result := jsondiff.Compare(gotJson, wantJson, &opts)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Expecting equal, diff: %s", result)
	}
}

func getPetJson() []byte {
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