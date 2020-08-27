package main

import (
	"bytes"
	"fmt"
	"net/http"

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
	methods := []resource.Method{}
	ct := resource.NewHTTPContentTypeSelector(resource.Response{415, resource.Response{415, "content-type not supported"}})
	ct.Add("application/json", encdec.JSONEncoderDecoder{}, true)

	//Resource
	res := resource.Resource{"pets", "/pets", []resource.Method{}}

	//Methods
	//Method Get (By id)
	failResponse := resource.Response{404, resource.Response{404, "Pet not found"}}
	successResponse := resource.Response{200, nil}
	mo := resource.NewMethodOperation(petstore.Pet{}, petstore.PetGetOperation{store}, successResponse, failResponse, getID, true)
	methodGet := resource.NewMethod(http.MethodGet, mo, ct)
	res.Methods = append(methods, methodGet)

	//Method Post (Create)
	failResponse = resource.Response{400, resource.Response{400, "Error creating pet"}}
	successResponse = resource.Response{201, nil}
	mo = resource.NewMethodOperation(petstore.Pet{}, petstore.PetCreateOperation{store}, successResponse, failResponse, getID, true)
	methodPost := resource.NewMethod(http.MethodPost, mo, ct)
	res.Methods = append(methods, methodPost)

	encoder := encdec.JSONEncoderDecoder{}
	//fmt.Println(res)
	buf := bytes.NewBufferString("")
	err := encoder.Encode(buf, res)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(buf)
}
