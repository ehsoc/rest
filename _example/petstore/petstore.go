package petstore

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"sync"

	"github.com/ehsoc/resource"
	"github.com/ehsoc/resource/encdec"
)

// Pet pet
//
// swagger:model Pet
type Pet struct {

	// id
	ID int64 `json:"id,omitempty" xml:"id,omitempty"`

	// name
	// Required: true
	Name string `json:"name" xml:"name"`

	// photo urls
	// Required: true
	PhotoUrls []string `json:"photoUrls" xml:"photoUrl"`

	// pet status in the store
	// Enum: [available pending sold]
	Status string `json:"status,omitempty" xml:"status,omitempty"`
}

type PetStore struct {
	Pets    map[int64]Pet
	idCount int64
	mutex   sync.Mutex
}

type PetGetOperation struct {
	Store PetStore
}

func (p PetGetOperation) Execute(id string, query url.Values, entityBody io.Reader, decoder encdec.Decoder) (interface{}, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	if pet, ok := p.Store.Pets[int64(idInt)]; ok {
		return pet, nil
	}
	return nil, errors.New("Pet not found.")
}

type PetCreateOperation struct {
	Store PetStore
}

func (p PetCreateOperation) Execute(id string, query url.Values, entityBody io.Reader, decoder encdec.Decoder) (interface{}, error) {
	pet := Pet{}
	err := decoder.Decode(entityBody, &pet)
	if err != nil {
		return nil, err
	}
	p.Store.mutex.Lock()
	p.Store.idCount++
	pet.ID = p.Store.idCount
	p.Store.Pets[pet.ID] = pet
	p.Store.mutex.Unlock()
	return pet, nil
}

func generatePetStore() resource.RestAPI {
	getIdFunc := func(r *http.Request) string {
		return "id"
	}
	api := resource.RestAPI{}
	api.BasePath = "/v2"
	api.Host = "localhost"
	pets, _ := resource.NewResource("/pet")
	contentTypes := resource.NewHTTPContentTypeSelector(resource.Response{})
	contentTypes.Add("application/json", encdec.JSONEncoderDecoder{}, true)
	contentTypes.Add("application/xml", encdec.JSONEncoderDecoder{}, false)
	//POST
	createMethodOperation := resource.NewMethodOperation(nil, resource.Response{201, nil, ""}, resource.Response{400, nil, ""}, true, true)
	createPetMethod := resource.NewMethod(http.MethodPost, createMethodOperation, contentTypes)
	createPetMethod.Summary = "Add a new pet to the store"
	createPetMethod.RequestBody = resource.RequestBody{"Pet object that needs to be added to the store", Pet{}}
	pets.AddMethod(createPetMethod)
	//New Resource with URIParam Resource GET By ID {petId} = /pet/{petId}
	eContentTypes := resource.NewHTTPContentTypeSelector(resource.Response{})
	eContentTypes.AddEncoder("application/json", encdec.JSONEncoderDecoder{}, true)
	eContentTypes.AddEncoder("application/xml", encdec.JSONEncoderDecoder{}, false)
	petIdResource, _ := resource.NewResourceWithURIParam("{petId}", getIdFunc, "", reflect.Int64)
	getByIdMethodOperation := resource.NewMethodOperation(nil, resource.Response{200, Pet{}, ""}, resource.Response{404, nil, ""}, true, false)
	getByIdPetMethod := resource.NewMethod(http.MethodGet, getByIdMethodOperation, eContentTypes)
	getByIdPetMethod.Summary = "Find pet by ID"
	getByIdPetMethod.Description = "Returns a single pet"
	petIdResource.GetURIParam().WithDescription("ID of pet to return")
	getByIdPetMethod.AddParameter(*petIdResource.GetURIParam())
	petIdResource.AddMethod(getByIdPetMethod)

	pets.Resources = append(pets.Resources, &petIdResource)
	//Delete
	deleteByIdMethodOperation := resource.NewMethodOperation(nil, resource.Response{200, nil, ""}, resource.Response{404, nil, ""}, false, false)
	deleteByIdMethod := resource.NewMethod(http.MethodDelete, deleteByIdMethodOperation, eContentTypes)
	deleteByIdMethod.Summary = "Deletes a pet"
	deleteByIdMethod.AddParameter(*petIdResource.GetURIParam().WithDescription("Pet id to delete"))
	apiKeyParam := resource.NewHeaderParameter("api_key", reflect.String, nil).AsOptional()
	deleteByIdMethod.AddParameter(*apiKeyParam)
	petIdResource.AddMethod(deleteByIdMethod)
	//Upload image resource under URIParameter Resource
	uploadImageResource, _ := resource.NewResource("uploadImage")
	uploadImageMethodOperation := resource.NewMethodOperation(nil, resource.Response{200, ApiResponse{}, "successful operation"}, resource.Response{}, false, false)
	eContentType := resource.NewHTTPContentTypeSelector(resource.Response{})
	eContentType.AddEncoder("application/json", encdec.JSONEncoderDecoder{}, true)
	eContentType.AddDecoder("multipart/form-data", encdec.XMLEncoderDecoder{}, true)
	uploadImageMethod := resource.NewMethod(http.MethodPost, uploadImageMethodOperation, eContentType)
	uploadImageMethod.Summary = "uploads an image"
	uploadImageMethod.AddParameter(*petIdResource.GetURIParam().WithDescription("ID of pet to update"))
	uploadImageMethod.AddParameter(*resource.NewFormDataParameter("additionalMetadata", reflect.String, encdec.JSONDecoder{}).WithDescription("Additional data to pass to server"))
	uploadImageMethod.AddParameter(*resource.NewFileParameter("file").WithDescription("file to upload"))
	uploadImageResource.AddMethod(uploadImageMethod)
	petIdResource.Resources = append(petIdResource.Resources, &uploadImageResource)

	api.Resources = append(api.Resources, &pets)
	return api
}
