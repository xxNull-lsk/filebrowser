package bolt

import (
	"fmt"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/favorite"
)

type favoriteBackend struct {
	db *storm.DB
}

func (s favoriteBackend) FindByUserID(id uint) ([]*favorite.Favorite, error) {
	var v []*favorite.Favorite
	err := s.db.Select(q.Eq("UserID", id)).Find(&v)
	if err == storm.ErrNotFound {
		return v, errors.ErrNotExist
	}

	return v, err
}

func (s favoriteBackend) Delete(hash string) error {
	err := s.db.DeleteStruct(&favorite.Favorite{Hash: hash})
	if err == storm.ErrNotFound {
		return nil
	}
	return err
}

func (s favoriteBackend) Get(path string, id uint) (*favorite.Favorite, error) {
	var v []*favorite.Favorite

	fmt.Printf("favoriteBackend: %s\n", path)
	err := s.db.Select(q.Eq("Path", path), q.Eq("UserID", id)).Find(v)
	if err == storm.ErrNotFound || len(v) == 0 {
		return nil, err
	}
	return v[0], err
}

func (s favoriteBackend) Save(f *favorite.Favorite) error {
	return s.db.Save(f)
}
