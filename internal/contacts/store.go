// internal/contacts/store.go
package contacts

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"sync"
)

var ErrNotFound = errors.New("contact introuvable")

type ContactStore interface {
	Add(c Contact) (Contact, error)
	Get(id int) (Contact, error)
	All() []Contact
	Update(id int, upd Contact) error
	Delete(id int) error
}

type Store struct {
	mu       sync.Mutex
	items    map[int]Contact
	nextID   int
	jsonPath string
}

// Verify that Store implements ContactStore
var _ ContactStore = (*Store)(nil)

func NewStore(jsonPath string) *Store {
	s := &Store{
		items:    make(map[int]Contact),
		nextID:   1,
		jsonPath: jsonPath,
	}
	_ = s.load()
	return s
}

func (s *Store) Add(c Contact) (Contact, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if c.ID == 0 {
		c.ID = s.nextID
		s.nextID++
	} else {
		if _, ok := s.items[c.ID]; ok {
			return Contact{}, errors.New("ID déjà utilisé: " + strconv.Itoa(c.ID))
		}
		if c.ID >= s.nextID {
			s.nextID = c.ID + 1
		}
	}

	s.items[c.ID] = c
	_ = s.save()
	return c, nil
}

func (s *Store) Get(id int) (Contact, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if c, ok := s.items[id]; ok {
		return c, nil
	}
	return Contact{}, ErrNotFound
}

func (s *Store) All() []Contact {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Contact, 0, len(s.items))
	for _, c := range s.items {
		out = append(out, c)
	}
	return out
}

func (s *Store) Update(id int, upd Contact) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}
	upd.ID = id
	s.items[id] = upd
	return s.save()
}

func (s *Store) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}
	delete(s.items, id)
	return s.save()
}

func (s *Store) load() error {
	f, err := os.Open(s.jsonPath)
	if err != nil {
		return err
	}
	defer f.Close()
	var arr []Contact
	if err := json.NewDecoder(f).Decode(&arr); err != nil {
		return err
	}
	maxID := 0
	for _, c := range arr {
		s.items[c.ID] = c
		if c.ID > maxID {
			maxID = c.ID
		}
	}
	s.nextID = maxID + 1
	return nil
}

func (s *Store) save() error {
	tmp := make([]Contact, 0, len(s.items))
	for _, c := range s.items {
		tmp = append(tmp, c)
	}
	f, err := os.Create(s.jsonPath)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(tmp)
}
