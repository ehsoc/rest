package main

import (
	"net/http"
	"os"

	"github.com/ehsoc/resource"
	"github.com/ehsoc/resource/_example/petstore"
	"github.com/ehsoc/resource/encdec"
	"github.com/go-chi/chi"
)

type ResponseBody struct {
	Code    int
	Message string
}

func getID(r *http.Request) string {
	return chi.URLParam(r, "petId")
}

func main() {
	store := petstore.PetStore{}
	store[1] = petstore.Pet{Name: "Dog"}
	methods := []resource.Method{}
	ct := resource.NewHTTPContentTypeSelector(resource.Response{415, Response{415, "content-type not supported"}})
	ct.Add("application/json", encdec.JSONEncoderDecoder{}, true)
	//Methods
	failResponse := resource.Response{404, Response{404, "Pet not found"}}
	successResponse := resource.Response{200, nil}
	mo := resource.NewMethodOperation(petstore.Pet{}, petstore.PetGetOperation{store}, successResponse, failResponse, getId, true)
	methodGet := resource.NewMethod(http.MethodGet, mo, ct)
	resource := resource.Resource{"pets", "/pets"}
	resource.Methods = append(methods, method)
	encdec.JSONEncoderDecoder.Encode(os.Stdout, resource)
}
