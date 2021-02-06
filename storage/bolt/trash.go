package bolt

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"

	"github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/trash"
)

type trashBackend struct {
	db *storm.DB
}

func (s trashBackend) All() ([]*trash.Trash, error) {
	var v []*trash.Trash
	err := s.db.All(&v)
	if err == storm.ErrNotFound {
		return v, errors.ErrNotExist
	}

	return v, err
}

func (s trashBackend) FindByUserID(id uint) ([]*trash.Trash, error) {
	var v []*trash.Trash
	err := s.db.Select(q.Eq("UserID", id)).Find(&v)
	if err == storm.ErrNotFound {
		return v, errors.ErrNotExist
	}

	return v, err
}

func (s trashBackend) GetByHash(hash string) (*trash.Trash, error) {
	var v trash.Trash
	err := s.db.One("Hash", hash, &v)
	if err == storm.ErrNotFound {
		return nil, errors.ErrNotExist
	}

	return &v, err
}

func (s trashBackend) Gets(path string, id uint) ([]*trash.Trash, error) {
	var v []*trash.Trash
	err := s.db.Select(q.Eq("Path", path), q.Eq("UserID", id)).Find(&v)
	if err == storm.ErrNotFound {
		return v, errors.ErrNotExist
	}

	return v, err
}

func (s trashBackend) Save(l *trash.Trash) error {
	return s.db.Save(l)
}

func (s trashBackend) Delete(hash string) error {
	err := s.db.DeleteStruct(&trash.Trash{Hash: hash})
	if err == storm.ErrNotFound {
		return nil
	}
	return err
}
