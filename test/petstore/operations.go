package petstore

import (
	"encoding/xml"
	"log"
	"strconv"

	"github.com/ehsoc/resource"
)

var PetStore = NewStore()

func operationCreate(i resource.Input) (interface{}, error) {
	pet := Pet{}
	body, _ := i.GetBody()
	err := i.BodyDecoder.Decode(body, &pet)
	if err != nil {
		return nil, err
	}
	return PetStore.Create(pet)
}

func operationUpdate(i resource.Input) (interface{}, error) {
	pet := Pet{}
	body, _ := i.GetBody()
	err := i.BodyDecoder.Decode(body, &pet)
	if err != nil {
		log.Println("error updating pet: ", err)
		return nil, err
	}
	pet, err = PetStore.Update(strconv.FormatInt(pet.ID, 10), pet)
	if err != nil {
		return pet, err
	}
	return pet, nil
}

func operationGetPetById(i resource.Input) (interface{}, error) {
	petId, _ := i.GetURIParam("petId")
	return PetStore.Get(petId)
}

func operationDeletePet(i resource.Input) (interface{}, error) {
	petId, _ := i.GetURIParam("petId")
	log.Println("Deleting pet id:", petId)
	return nil, PetStore.Delete(petId)
}

func operationUploadImage(i resource.Input) (interface{}, error) {
	petId, _ := i.GetURIParam("petId")
	log.Println("Uploading image pet id:", petId)
	fb, _, _ := i.GetFormFile("file")
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

func operationFindByStatus(i resource.Input) (interface{}, error) {
	status, err := i.GetQuery("status")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("searching by status: ")
	for _, s := range status {
		log.Print(s, " ")
	}
	petsList, err := PetStore.List()
	if err != nil {
		return nil, err
	}
	//If the encoder is XML, we want to wrap it with <pets>
	if i.Request.Context().Value(resource.ContentTypeContextKey("encoder")) == "application/xml" {
		return XmlPetsWrapper{Pets: petsList}, nil
	}
	return petsList, nil
}
