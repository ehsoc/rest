package petstore

import (
	"encoding/xml"
	"log"
	"net/http"

	"github.com/ehsoc/resource"
	"github.com/ehsoc/resource/encdec"
	"github.com/ehsoc/resource/httputil"
	"github.com/go-chi/chi"
)

var PetStore = NewStore()

func operationCreate(r *http.Request, decoder encdec.Decoder) (interface{}, error) {
	pet := Pet{}
	err := decoder.Decode(r.Body, &pet)
	if err != nil {
		return nil, err
	}
	return PetStore.Create(pet)
}

func operationGetPetById(r *http.Request, decoder encdec.Decoder) (interface{}, error) {
	petId := chi.URLParam(r, "petId")
	return PetStore.Get(petId)
}

func operationDeletePet(r *http.Request, decoder encdec.Decoder) (interface{}, error) {
	petId := chi.URLParam(r, "petId")
	log.Println("Deleting pet id:", petId)
	return nil, PetStore.Delete(petId)
}

func operationUploadImage(r *http.Request, decoder encdec.Decoder) (interface{}, error) {
	petId := chi.URLParam(r, "petId")
	log.Println("Uploading image pet id:", petId)
	fb, _, _ := httputil.GetFormFile(r, "file")
	err := PetStore.UploadPhoto(petId, fb)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

type XmlPetsWrapper struct {
	XMLName xml.Name `xml:"pets"`
	Pets    []Pet    `xml:"Pet"`
}

func operationFindByStatus(r *http.Request, decoder encdec.Decoder) (interface{}, error) {
	if status, ok := r.URL.Query()["status"]; ok {
		log.Println("searching by status: ")
		for _, s := range status {
			log.Print(s, " ")
		}
	}
	petsList, err := PetStore.List()
	if err != nil {
		return nil, err
	}
	//If the encoder is XML, we want to wrap it with <pets>
	if r.Context().Value(resource.ContentTypeContextKey("encoder")) == "application/xml" {
		return XmlPetsWrapper{Pets: petsList}, nil
	}
	return petsList, nil
}
