package contacts

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type GormStore struct {
	db *gorm.DB
}

// Vérifie que GormStore implémente ContactStore
var _ ContactStore = (*GormStore)(nil)

func NewGormStore(db *gorm.DB) (*GormStore, error) {
	// Auto-migration: crée/maj le schéma si nécessaire
	if err := db.AutoMigrate(&Contact{}); err != nil {
		return nil, fmt.Errorf("automigrate: %w", err)
	}
	return &GormStore{db: db}, nil
}

func (s *GormStore) Add(c Contact) (Contact, error) {
	if !c.IsValid() {
		return Contact{}, errors.New("nom et email sont obligatoires")
	}
	// si c.ID == 0, GORM auto-incrémente
	if err := s.db.Create(&c).Error; err != nil {
		return Contact{}, err
	}
	return c, nil
}

func (s *GormStore) Get(id int) (Contact, error) {
	var out Contact
	err := s.db.First(&out, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Contact{}, ErrNotFound
		}
		return Contact{}, err
	}
	return out, nil
}

func (s *GormStore) All() []Contact {
	var list []Contact
	_ = s.db.Order("id ASC").Find(&list).Error
	return list
}

func (s *GormStore) Update(id int, upd Contact) error {
	var existing Contact
	if err := s.db.First(&existing, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	if upd.Nom != "" {
		existing.Nom = upd.Nom
	}
	if upd.Email != "" {
		existing.Email = upd.Email
	}
	if !existing.IsValid() {
		return errors.New("nom et email ne peuvent pas être vides")
	}

	if err := s.db.Save(&existing).Error; err != nil {
		// Optionnel : message plus clair si email dupliqué
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("email déjà utilisé: %s", existing.Email)
		}
		return err
	}
	return nil
}

func (s *GormStore) Delete(id int) error {
	res := s.db.Delete(&Contact{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
