package petstore

import (
	"bytes"
	"io"
	"io/ioutil"
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
	createMethodOperation := resource.NewMethodOperation(resource.OperationFunc(operationCreate), resource.NewResponse(201)).WithFailResponse(resource.NewResponse(400))
	pets.Post(createMethodOperation, contentTypes).
		WithRequestBody("Pet object that needs to be added to the store", Pet{}).
		WithSummary("Add a new pet to the store")
	//PUT
	updateMethodOperation := resource.NewMethodOperation(resource.OperationFunc(operationUpdate), resource.NewResponse(200)).WithFailResponse(resource.NewResponse(404).WithDescription("Pet not found"))
	pets.Put(updateMethodOperation, contentTypes).
		WithRequestBody("Pet object that needs to be added to the store", Pet{}).
		WithSummary("Update an existing pet").
		WithValidation(resource.ValidatorFunc(func(input resource.Input) error {
			pet := Pet{}
			body, _ := input.GetBody()
			respBody := new(bytes.Buffer)
			cBody := io.TeeReader(body, respBody)
			err := input.BodyDecoder.Decode(cBody, &pet)
			if err != nil {
				return err
			}
			input.Request.Body = ioutil.NopCloser(respBody)
			return nil
		}),
			resource.NewResponse(400).WithDescription("Invalid ID supplied"))

	//Uri Parameters declaration, so it is available to all anonymous resources functions
	petIdURIParam := resource.NewURIParameter("petId", reflect.Int64).WithDescription("ID of pet to return").WithExample(1)
	//SubResource
	//New Resource with URIParam Resource GET By ID {petId} = /pet/{petId}
	pets.Resource("{petId}", func(r *resource.Resource) {
		ct := resource.NewHTTPContentTypeSelector()
		ct.AddEncoder("application/json", encdec.JSONEncoder{}, true)
		ct.AddEncoder("application/xml", encdec.XMLEncoder{}, false)
		getByIdMethodOperation := resource.NewMethodOperation(resource.OperationFunc(operationGetPetById), resource.NewResponse(200).WithOperationResultBody(Pet{})).WithFailResponse(notFoundResponse)
		r.Get(getByIdMethodOperation, ct).
			WithSummary("Find pet by ID").
			WithDescription("Returns a single pet").
			WithParameter(petIdURIParam)
		//Delete
		deleteByIdMethodOperation := resource.NewMethodOperation(resource.OperationFunc(operationDeletePet), resource.NewResponse(200)).WithFailResponse(notFoundResponse)
		r.Delete(deleteByIdMethodOperation, ct).
			WithSummary("Deletes a pet").
			WithParameter(
				petIdURIParam.WithDescription("Pet id to delete").
					WithValidation(resource.ValidatorFunc(func(i resource.Input) error {
						petId, _ := i.GetURIParam("petId")
						_, err := getInt64Id(petId)
						if err != nil {
							return err
						}
						return nil
					}), resource.NewResponse(400).WithDescription("Invalid ID supplied"))).
			WithParameter(resource.NewHeaderParameter("api_key", reflect.String))
		r.Resource("uploadImage", func(r *resource.Resource) {
			//Upload image resource under URIParameter Resource
			uploadImageMethodOperation := resource.NewMethodOperation(resource.OperationFunc(operationUploadImage), resource.NewResponse(200).WithBody(ApiResponse{200, "OK", "image created"}).WithDescription("successful operation"))
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
		findByStatusMethodOperation := resource.NewMethodOperation(resource.OperationFunc(operationFindByStatus), resource.NewResponse(200).WithOperationResultBody([]Pet{}).WithDescription("successful operation")).WithFailResponse(resource.NewResponse(400).WithDescription("Invalid status value"))
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
