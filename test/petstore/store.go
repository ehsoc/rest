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

func getInt64Id(stringId string) (int64, error) {
	id, err := strconv.Atoi(stringId)
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}

func (s *Store) Get(petId string) (Pet, error) {
	id, err := getInt64Id(petId)
	if err != nil {
		return Pet{}, err
	}
	if pet, ok := s.store[id]; ok {
		return pet, nil
	}
	return Pet{}, errors.New("Pet not found.")
}

func (s *Store) Create(pet Pet) (Pet, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.idCount++
	pet.ID = s.idCount
	s.store[pet.ID] = pet
	return pet, nil
}

func (s *Store) Delete(petId string) error {
	id, err := getInt64Id(petId)
	if err != nil {
		return err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.store[id]; !ok {
		return errors.New("pet not found")
	}
	delete(s.store, id)
	return nil
}

func (s *Store) List() ([]Pet, error) {
	list := []Pet{}
	for _, pet := range s.store {
		list = append(list, pet)
	}
	return list, nil
}
