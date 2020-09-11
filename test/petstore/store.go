package petstore

import (
	"errors"
	"strconv"
	"sync"
)

//The objective of this package is to provide a CRUD store a non-persistent store in memory for integration tests.
type Store struct {
	store   map[int64]Pet
	idCount int64
	mutex   sync.Mutex
}

func NewStore() Store {
	store := Store{}
	store.store = make(map[int64]Pet)
	return store
}

func (s *Store) Get(petId string) (Pet, error) {
	id, err := strconv.Atoi(petId)
	if err != nil {
		return Pet{}, err
	}
	if pet, ok := s.store[int64(id)]; ok {
		return pet, nil
	}
	return Pet{}, errors.New("Pet not found.")
}

func (s *Store) Create(pet Pet) (Pet, error) {
	s.mutex.Lock()
	s.idCount++
	pet.ID = s.idCount
	s.store[pet.ID] = pet
	s.mutex.Unlock()
	return pet, nil
}

func (s *Store) List() ([]Pet, error) {
	list := []Pet{}
	for _, pet := range s.store {
		list = append(list, pet)
	}
	return list, nil
}
