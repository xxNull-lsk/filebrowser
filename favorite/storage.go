package favorite

// StorageBackend is the interface to implement for a favorite storage.
type StorageBackend interface {
	FindByUserID(id uint) ([]*Favorite, error)
	Get(path string, id uint) (*Favorite, error)
	Save(f *Favorite) error
	Delete(path string) error
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

// Get wraps a StorageBackend.Get
func (s *Storage) Get(path string, id uint) (*Favorite, error) {
	return s.back.Get(path, id)
}

// Delete wraps a StorageBackend.Delete
func (s *Storage) Delete(hash string) error {
	return s.back.Delete(hash)
}
