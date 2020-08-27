package petstore

import (
	"errors"
	"io"
	"net/url"
	"strconv"
	"sync"

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
