package petstore

import (
	"reflect"

	"github.com/ehsoc/resource"
	"github.com/ehsoc/resource/encdec"
)

type ApiResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type dummyBody struct{}

var unsupportedResponse = resource.NewResponse(415)
var notFoundResponse = resource.NewResponse(404)

func GeneratePetStore() resource.RestAPI {
	//Resource /pet
	pets := resource.NewResource("pet")
	contentTypes := resource.NewHTTPContentTypeSelector()
	contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
	contentTypes.Add("application/xml", encdec.XMLEncoderDecoder{}, false)
	//POST
	createMethodOperation := resource.NewMethodOperation(resource.OperationFunc(operationCreate), resource.NewResponse(201), resource.NewResponse(400), true)
	pets.Post(createMethodOperation, contentTypes).
		WithRequestBody("Pet object that needs to be added to the store", Pet{}).
		WithSummary("Add a new pet to the store")
	//PUT
	updateMethodOperation := resource.NewMethodOperation(resource.OperationFunc(operationUpdate), resource.NewResponse(200), resource.NewResponse(404).WithDescription("Pet not found"), true)
	pets.Put(updateMethodOperation, contentTypes).
		WithRequestBody("Pet object that needs to be added to the store", Pet{}).
		WithSummary("Update an existing pet").
		WithValidation(resource.ValidatorFunc(func(input resource.Input) error {
			//body, _ := input.GetBody()
			return nil
		}),
			resource.NewResponse(400).WithDescription("Invalid ID supplied"))

	//Uri Parameters declaration, so it is available to all anonymous resources functions
	petIdURIParam := resource.NewURIParameter("petId", reflect.Int64).WithDescription("ID of pet to return")
	//SubResource
	//New Resource with URIParam Resource GET By ID {petId} = /pet/{petId}
	pets.Resource("{petId}", func(r *resource.Resource) {
		ct := resource.NewHTTPContentTypeSelector()
		ct.AddEncoder("application/json", encdec.JSONEncoder{}, true)
		ct.AddEncoder("application/xml", encdec.XMLEncoder{}, false)
		getByIdMethodOperation := resource.NewMethodOperation(resource.OperationFunc(operationGetPetById), resource.NewResponse(200).WithBody(Pet{}), notFoundResponse, true)
		r.Get(getByIdMethodOperation, ct).
			WithSummary("Find pet by ID").
			WithDescription("Returns a single pet").
			WithParameter(petIdURIParam)
		//Delete
		deleteByIdMethodOperation := resource.NewMethodOperation(resource.OperationFunc(operationDeletePet), resource.NewResponse(200), notFoundResponse, false)
		r.Delete(deleteByIdMethodOperation, ct).
			WithSummary("Deletes a pet").
			WithParameter(petIdURIParam.WithDescription("Pet id to delete")).
			WithParameter(resource.NewHeaderParameter("api_key", reflect.String).AsOptional())
		r.Resource("uploadImage", func(r *resource.Resource) {
			//Upload image resource under URIParameter Resource
			uploadImageMethodOperation := resource.NewMethodOperation(resource.OperationFunc(operationUploadImage), resource.NewResponse(200).WithBody(ApiResponse{200, "OK", "image created"}).WithDescription("successful operation"), resource.NewResponse(500), false)
			ct := resource.NewHTTPContentTypeSelector()
			ct.AddEncoder("application/json", encdec.JSONEncoderDecoder{}, true)
			ct.AddDecoder("multipart/form-data", encdec.XMLEncoderDecoder{}, true)
			r.Post(uploadImageMethodOperation, ct).
				WithParameter(petIdURIParam.WithDescription("ID of pet to update")).
				WithParameter(resource.NewFormDataParameter("additionalMetadata", reflect.String, encdec.JSONDecoder{}).WithDescription("Additional data to pass to server")).
				WithParameter(resource.NewFileParameter("file").WithDescription("file to upload")).
				WithParameter(resource.NewFormDataParameter("jsonPetData", reflect.Struct, encdec.JSONDecoder{}).WithDescription("json format data").WithBody(Pet{})).
				WithSummary("uploads an image")
		})
	})
	//Resource /pet/findByStatus
	pets.Resource("findByStatus", func(r *resource.Resource) {
		ct := resource.NewHTTPContentTypeSelector()
		ct.AddEncoder("application/json", encdec.JSONEncoderDecoder{}, true)
		ct.AddEncoder("application/xml", encdec.XMLEncoderDecoder{}, false)
		//GET
		findByStatusMethodOperation := resource.NewMethodOperation(resource.OperationFunc(operationFindByStatus), resource.NewResponse(200).WithBody([]Pet{}).WithDescription("successful operation"), resource.NewResponse(400).WithDescription("Invalid status value"), true)
		statusParam := resource.NewQueryArrayParameter("status", []interface{}{"available", "pending", "sold"}).AsRequired().WithDescription("Status values that need to be considered for filter")
		statusParam.CollectionFormat = "multi"
		r.Get(findByStatusMethodOperation, ct).
			WithSummary("Finds Pets by status").
			WithDescription("Multiple status values can be provided with comma separated strings").
			WithParameter(statusParam)
	})

	api := resource.RestAPI{}
	api.BasePath = "/v2"
	api.Host = "localhost"
	api.AddResource(pets)
	return api
}
