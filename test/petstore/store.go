package petstore

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/afero"
)

//The objective of this package is to provide a CRUD store a non-persistent store in memory for integration tests.
type Store struct {
	store      map[int64]Pet
	idCount    int64
	mutex      sync.Mutex
	InMemoryFs afero.Fs
}

var ErrorPetNotFound = errors.New("Pet not found.")

func NewStore() *Store {
	store := Store{}
	store.store = make(map[int64]Pet)
	store.InMemoryFs = afero.NewMemMapFs()
	return &store
}

func getInt64Id(stringId string) (int64, error) {
	id, err := strconv.Atoi(stringId)
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}

func (s *Store) Get(petID string) (Pet, error) {
	id, err := getInt64Id(petID)
	if err != nil {
		return Pet{}, err
	}
	log.Printf("searching pet id: %d\n", id)
	if pet, ok := s.store[id]; ok {
		log.Printf("pet found id: %d\n", id)
		return pet, nil
	}
	return Pet{}, ErrorPetNotFound
}

func (s *Store) Create(pet Pet) (Pet, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.idCount++
	pet.ID = s.idCount
	s.store[pet.ID] = pet
	return pet, nil
}

func (s *Store) Update(petID string, pet Pet) (Pet, error) {
	id, err := getInt64Id(petID)
	if err != nil {
		return Pet{}, err
	}
	petFound, err := s.Get(petID)
	if err != nil {
		return petFound, err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.store[id] = pet
	return pet, nil
}

func (s *Store) Delete(petID string) error {
	id, err := getInt64Id(petID)
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

func (s *Store) UploadPhoto(petID string, fileContent []byte) error {
	id, err := getInt64Id(petID)
	if err != nil {
		return err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if pet, ok := s.store[id]; ok {
		url := fmt.Sprintf("files/%s%d", petID, time.Now().UnixNano())
		err := afero.WriteFile(s.InMemoryFs, url, fileContent, 0655)
		if err != nil {
			return err
		}
		pet.PhotoUrls = append(s.store[id].PhotoUrls, url)
		s.store[id] = pet
	}
	return nil
}
