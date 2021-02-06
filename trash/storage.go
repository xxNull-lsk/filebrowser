package trash

// StorageBackend is the interface to implement for a share storage.
type StorageBackend interface {
	All() ([]*Trash, error)
	FindByUserID(id uint) ([]*Trash, error)
	GetByHash(hash string) (*Trash, error)
	Gets(path string, id uint) ([]*Trash, error)
	Save(s *Trash) error
	Delete(hash string) error
}

// Storage is a storage.
type Storage struct {
	back StorageBackend
}

// NewStorage creates a share links storage from a backend.
func NewStorage(back StorageBackend) *Storage {
	return &Storage{back: back}
}

// All wraps a StorageBackend.All.
func (s *Storage) All() ([]*Trash, error) {
	trashs, err := s.back.All()

	if err != nil {
		return nil, err
	}

	return trashs, nil
}

// FindByUserID wraps a StorageBackend.FindByUserID.
func (s *Storage) FindByUserID(id uint) ([]*Trash, error) {
	trashs, err := s.back.FindByUserID(id)

	if err != nil {
		return nil, err
	}

	return trashs, nil
}

// GetByHash wraps a StorageBackend.GetByHash.
func (s *Storage) GetByHash(hash string) (*Trash, error) {
	trash, err := s.back.GetByHash(hash)
	if err != nil {
		return nil, err
	}

	return trash, nil
}

// Gets wraps a StorageBackend.Gets
func (s *Storage) Gets(path string, id uint) ([]*Trash, error) {
	trashs, err := s.back.Gets(path, id)

	if err != nil {
		return nil, err
	}

	return trashs, nil
}

// Save wraps a StorageBackend.Save
func (s *Storage) Save(l *Trash) error {
	return s.back.Save(l)
}

// Delete wraps a StorageBackend.Delete
func (s *Storage) Delete(hash string) error {
	return s.back.Delete(hash)
}
