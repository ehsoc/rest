package petstore

import (
	"io"
	"log"
	"net/url"

	"github.com/ehsoc/resource/encdec"
)

var PetStore = NewStore()

func operationCreate(body io.ReadCloser, parameters url.Values, decoder encdec.Decoder) (interface{}, error) {
	pet := Pet{}
	err := decoder.Decode(body, &pet)
	if err != nil {
		return nil, err
	}
	return PetStore.Create(pet)
}

func operationGetPetById(body io.ReadCloser, parameters url.Values, decoder encdec.Decoder) (interface{}, error) {
	log.Println("Searching pet id:", parameters.Get("petId"))
	return PetStore.Get(parameters.Get("petId"))
}

func operationDeletePet(body io.ReadCloser, parameters url.Values, decoder encdec.Decoder) (interface{}, error) {
	log.Println("Deleting pet id:", parameters.Get("petId"))
	return nil, PetStore.Delete(parameters.Get("petId"))
}

func operationUploadImage(body io.ReadCloser, parameters url.Values, decoder encdec.Decoder) (interface{}, error) {
	petId := parameters.Get("petId")
	log.Println("Uploading image pet id:", parameters.Get("petId"))
	fileString := parameters.Get("file")
	err := PetStore.UploadPhoto(petId, fileString)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
