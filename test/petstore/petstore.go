package petstore

import (
	"net/http"
	"reflect"

	"github.com/ehsoc/resource"
	"github.com/ehsoc/resource/encdec"
)

type ApiResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

var unsupportedResponse = resource.Response{415, nil, ""}

func NewPetStoreResource() resource.Resource {
	r, _ := resource.NewResource("/api/v2")
	r.Resource("/pet", func(r *resource.Resource) {
		contentTypes := resource.NewHTTPContentTypeSelector(unsupportedResponse)
		contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
		contentTypes.Add("application/xml", encdec.XMLEncoderDecoder{}, false)
		createMethodOperation := resource.NewMethodOperation(resource.OperationFunc(operationCreate), resource.Response{201, nil, ""}, resource.Response{400, nil, ""}, true)
		createPetMethod := resource.NewMethod(http.MethodPost, createMethodOperation, contentTypes)
		createPetMethod.Summary = "Add a new pet to the store"
		createPetMethod.RequestBody = resource.RequestBody{"Pet object that needs to be added to the store", Pet{}}
		r.AddMethod(createPetMethod)
	})
	return r
}

func GeneratePetStore() resource.RestAPI {
	api := resource.RestAPI{}
	api.BasePath = "/v2"
	api.Host = "localhost"
	//Resource /pet
	pets, _ := resource.NewResource("/pet")
	contentTypes := resource.NewHTTPContentTypeSelector(unsupportedResponse)
	contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
	contentTypes.Add("application/xml", encdec.XMLEncoderDecoder{}, false)
	//POST
	createMethodOperation := resource.NewMethodOperation(resource.OperationFunc(operationCreate), resource.Response{201, nil, ""}, resource.Response{400, nil, ""}, true)
	createPetMethod := resource.NewMethod(http.MethodPost, createMethodOperation, contentTypes)
	createPetMethod.Summary = "Add a new pet to the store"
	createPetMethod.RequestBody = resource.RequestBody{"Pet object that needs to be added to the store", Pet{}}
	pets.AddMethod(createPetMethod)
	//New Resource with URIParam Resource GET By ID {petId} = /pet/{petId}
	eContentTypes := resource.NewHTTPContentTypeSelector(unsupportedResponse)
	eContentTypes.AddEncoder("application/json", encdec.JSONEncoderDecoder{}, true)
	eContentTypes.AddEncoder("application/xml", encdec.XMLEncoderDecoder{}, false)
	petIdResource, _ := resource.NewResourceWithURIParam("/{petId}", "", reflect.Int64)
	getByIdMethodOperation := resource.NewMethodOperation(resource.OperationFunc(operationGetPetById), resource.Response{200, Pet{}, ""}, resource.Response{404, nil, ""}, true)
	getByIdPetMethod := resource.NewMethod(http.MethodGet, getByIdMethodOperation, eContentTypes)
	getByIdPetMethod.Summary = "Find pet by ID"
	getByIdPetMethod.Description = "Returns a single pet"
	petIdResource.GetURIParam().WithDescription("ID of pet to return")
	getByIdPetMethod.AddParameter(*petIdResource.GetURIParam())
	petIdResource.AddMethod(getByIdPetMethod)

	pets.Resources = append(pets.Resources, &petIdResource)
	//Delete
	deleteByIdMethodOperation := resource.NewMethodOperation(resource.OperationFunc(operationDeletePet), resource.Response{200, nil, ""}, resource.Response{404, nil, ""}, false)
	deleteByIdMethod := resource.NewMethod(http.MethodDelete, deleteByIdMethodOperation, eContentTypes)
	deleteByIdMethod.Summary = "Deletes a pet"
	deleteByIdMethod.AddParameter(*petIdResource.GetURIParam().WithDescription("Pet id to delete"))
	apiKeyParam := resource.NewHeaderParameter("api_key", reflect.String).AsOptional()
	deleteByIdMethod.AddParameter(*apiKeyParam)
	petIdResource.AddMethod(deleteByIdMethod)
	//Upload image resource under URIParameter Resource
	uploadImageResource, _ := resource.NewResource("/uploadImage")
	uploadImageMethodOperation := resource.NewMethodOperation(resource.OperationFunc(operationUploadImage), resource.Response{200, ApiResponse{200, "OK", "image created"}, "successful operation"}, resource.Response{500, nil, ""}, false)
	eContentType := resource.NewHTTPContentTypeSelector(unsupportedResponse)
	eContentType.AddEncoder("application/json", encdec.JSONEncoderDecoder{}, true)
	eContentType.AddDecoder("multipart/form-data", encdec.XMLEncoderDecoder{}, true)
	uploadImageMethod := resource.NewMethod(http.MethodPost, uploadImageMethodOperation, eContentType)
	uploadImageMethod.Summary = "uploads an image"
	uploadImageMethod.AddParameter(*petIdResource.GetURIParam().WithDescription("ID of pet to update"))
	uploadImageMethod.AddParameter(*resource.NewFormDataParameter("additionalMetadata", reflect.String, encdec.JSONDecoder{}).WithDescription("Additional data to pass to server"))
	uploadImageMethod.AddParameter(*resource.NewFileParameter("file").WithDescription("file to upload"))
	uploadImageMethod.AddParameter(*resource.NewFormDataParameter("jsonPetData", reflect.Struct, encdec.JSONDecoder{}).WithDescription("json format data").WithBody(Pet{}))
	uploadImageResource.AddMethod(uploadImageMethod)
	petIdResource.Resources = append(petIdResource.Resources, &uploadImageResource)
	//Resource /pet/findByStatus
	petsFindRes, _ := resource.NewResource("/pet/findByStatus")
	contentTypesF := resource.NewHTTPContentTypeSelector(unsupportedResponse)
	contentTypesF.AddEncoder("application/json", encdec.JSONEncoderDecoder{}, true)
	contentTypesF.AddEncoder("application/xml", encdec.XMLEncoderDecoder{}, false)
	//GET
	findByStatusMethodOperation := resource.NewMethodOperation(resource.OperationFunc(operationFindByStatus), resource.Response{200, []Pet{}, "successful operation"}, resource.Response{400, nil, "Invalid status value"}, true)
	findByStatusPetMethod := resource.NewMethod(http.MethodGet, findByStatusMethodOperation, contentTypesF)
	findByStatusPetMethod.Summary = "Finds Pets by status"
	findByStatusPetMethod.Description = "Multiple status values can be provided with comma separated strings"
	statusParam := resource.NewQueryArrayParameter("status", []interface{}{"available", "pending", "sold"}).AsRequired().WithDescription("Status values that need to be considered for filter")
	statusParam.CollectionFormat = "multi"
	findByStatusPetMethod.AddParameter(*statusParam)
	petsFindRes.AddMethod(findByStatusPetMethod)
	api.Resources = append(api.Resources, &petsFindRes)
	api.Resources = append(api.Resources, &pets)
	return api
}
