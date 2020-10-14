package petstore

import (
	"encoding/xml"
	"log"
	"strconv"

	"github.com/ehsoc/resource"
)

var PetStore = NewStore()

func operationCreate(i resource.Input) (interface{}, bool, error) {
	pet := Pet{}
	body, _ := i.GetBody()
	err := i.BodyDecoder.Decode(body, &pet)
	if err != nil {
		return nil, false, err
	}
	pet, err = PetStore.Create(pet)
	if err != nil {
		return pet, false, err
	}
	return pet, true, nil
}

func operationUpdate(i resource.Input) (interface{}, bool, error) {
	pet := Pet{}
	body, _ := i.GetBody()
	err := i.BodyDecoder.Decode(body, &pet)
	if err != nil {
		log.Println("error updating pet: ", err)
		return nil, false, err
	}
	pet, err = PetStore.Update(strconv.FormatInt(pet.ID, 10), pet)
	if err != nil {
		return pet, false, err
	}
	return pet, true, nil
}

func operationGetPetById(i resource.Input) (interface{}, bool, error) {
	petId, err := i.GetURIParam("petId")
	if err != nil {
		log.Fatal(err)
		return nil, false, err
	}
	pet, err := PetStore.Get(petId)
	if err != nil {
		if err == ErrorPetNotFound {
			//not found but is not an error
			return pet, false, nil
		}
		return pet, false, err
	}
	return pet, true, nil
}

func operationDeletePet(i resource.Input) (interface{}, bool, error) {
	petId, _ := i.GetURIParam("petId")
	log.Println("Deleting pet id:", petId)
	err := PetStore.Delete(petId)
	if err != nil {
		return nil, false, err
	}
	return nil, true, nil
}

func operationUploadImage(i resource.Input) (interface{}, bool, error) {
	petId, _ := i.GetURIParam("petId")
	log.Println("Uploading image pet id:", petId)
	fb, _, _ := i.GetFormFile("file")
	err := PetStore.UploadPhoto(petId, fb)
	if err != nil {
		return nil, false, err
	}
	return nil, true, nil
}

type XmlPetsWrapper struct {
	XMLName xml.Name `xml:"pets"`
	Pets    []Pet    `xml:"Pet"`
}

func operationFindByStatus(i resource.Input) (interface{}, bool, error) {
	status, err := i.GetQuery("status")
	if err != nil {
		return nil, false, err
	}
	log.Println("searching by status: ")
	for _, s := range status {
		log.Print(s, " ")
	}
	petsList, err := PetStore.List()
	if err != nil {
		return nil, false, err
	}
	//If the encoder is XML, we want to wrap it with <pets>
	if i.Request.Context().Value(resource.ContentTypeContextKey("encoder")) == "application/xml" {
		return XmlPetsWrapper{Pets: petsList}, true, nil
	}
	return petsList, true, nil
}
