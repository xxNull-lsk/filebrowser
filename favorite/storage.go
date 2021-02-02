package favorite

import (
	"time"

	"github.com/filebrowser/filebrowser/v2/errors"
)

// StorageBackend is the interface to implement for a share storage.
type StorageBackend interface {
	FindByUserID(id uint) ([]*Favorite, error)
	Has(path string, id uint) error
	Save(f *Favorite) error
	Delete(path string, id uint) error
}

// Storage is a storage.
type Storage struct {
	back StorageBackend
}

// NewStorage creates a favorites storage from a backend.
func NewStorage(back StorageBackend) *Storage {
	return &Storage{back: back}
}

// FindByUserID wraps a StorageBackend.FindByUserID.
func (s *Storage) FindByUserID(id uint) ([]*Favorite, error) {
	links, err := s.back.FindByUserID(id)

	if err != nil {
		return nil, err
	}

	return links, nil
}

// Save wraps a StorageBackend.Save
func (s *Storage) Save(f *Favorite) error {
	return s.back.Save(f)
}

// Has wraps a StorageBackend.Has
func (s *Storage) Has(path string, id uint) error {
	return s.back.Has(path, id)
}

// Delete wraps a StorageBackend.Delete
func (s *Storage) Delete(path string, id uint) error {
	return s.back.Delete(path, id)
}
