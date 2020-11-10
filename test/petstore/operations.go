package petstore

import (
	"encoding/xml"
	"log"
	"strconv"

	"github.com/ehsoc/rest"
)

var PetStore = NewStore()

func operationCreate(i rest.Input) (interface{}, bool, error) {
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

func operationUpdate(i rest.Input) (interface{}, bool, error) {
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

func operationGetPetByID(i rest.Input) (interface{}, bool, error) {
	petID, err := i.GetURIParam("petId")
	if err != nil {
		log.Fatal(err)
		return nil, false, err
	}
	pet, err := PetStore.Get(petID)
	if err != nil {
		if err == ErrorPetNotFound {
			// not found but is not an error
			return pet, false, nil
		}
		return pet, false, err
	}
	return pet, true, nil
}

func operationDeletePet(i rest.Input) (interface{}, bool, error) {
	petID, _ := i.GetURIParam("petId")
	log.Println("Deleting pet id:", petID)
	err := PetStore.Delete(petID)
	if err != nil {
		return nil, false, err
	}
	return nil, true, nil
}

func operationUploadImage(i rest.Input) (interface{}, bool, error) {
	petID, _ := i.GetURIParam("petId")
	log.Println("Uploading image pet id:", petID)
	fb, _, _ := i.GetFormFile("file")
	err := PetStore.UploadPhoto(petID, fb)
	if err != nil {
		return nil, false, err
	}
	return nil, true, nil
}

type XMLPetsWrapper struct {
	XMLName xml.Name `xml:"pets"`
	Pets    []Pet    `xml:"Pet"`
}

func operationFindByStatus(i rest.Input) (interface{}, bool, error) {
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
	// If the encoder is XML, we want to wrap it with <pets>
	if i.Request.Context().Value(rest.ContentTypeContextKey("encoder")) == "application/xml" {
		return XMLPetsWrapper{Pets: petsList}, true, nil
	}
	return petsList, true, nil
}
