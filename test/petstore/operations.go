package petstore

import (
	"io"
	"log"
	"net/url"

	"github.com/ehsoc/resource/encdec"
)

var store = NewStore()

func operationCreate(body io.ReadCloser, parameters url.Values, decoder encdec.Decoder) (interface{}, error) {
	pet := Pet{}
	err := decoder.Decode(body, &pet)
	if err != nil {
		return nil, err
	}
	return store.Create(pet)
}

func operationGetPetById(body io.ReadCloser, parameters url.Values, decoder encdec.Decoder) (interface{}, error) {
	log.Println("Searching pet id:", parameters.Get("petId"))
	return store.Get(parameters.Get("petId"))
}
