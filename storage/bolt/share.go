package bolt

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"

	"github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/share"
)

type shareBackend struct {
	db *storm.DB
}

func (s shareBackend) All() ([]*share.Link, error) {
	var v []*share.Link
	err := s.db.All(&v)
	if err == storm.ErrNotFound {
		return v, errors.ErrNotExist
	}

	return v, err
}

func (s shareBackend) FindByUserID(id uint) ([]*share.Link, error) {
	var v []*share.Link
	err := s.db.Select(q.Eq("UserID", id)).Find(&v)
	if err == storm.ErrNotFound {
		return v, errors.ErrNotExist
	}

	return v, err
}

func (s shareBackend) GetByHash(hash string) (*share.Link, error) {
	var v share.Link
	err := s.db.One("Hash", hash, &v)
	if err == storm.ErrNotFound {
		return nil, errors.ErrNotExist
	}

	return &v, err
}

func (s shareBackend) GetPermanent(path string, id uint) (*share.Link, error) {
	var v share.Link
	err := s.db.Select(q.Eq("Path", path), q.Eq("Expire", 0), q.Eq("UserID", id)).First(&v)
	if err == storm.ErrNotFound {
		return nil, errors.ErrNotExist
	}

	return &v, err
}

func (s shareBackend) Gets(path string, id uint) ([]*share.Link, error) {
	var v []*share.Link
	err := s.db.Select(q.Eq("Path", path), q.Eq("UserID", id)).Find(&v)
	if err == storm.ErrNotFound {
		return v, errors.ErrNotExist
	}

	return v, err
}

func (s shareBackend) Save(l *share.Link) error {
	return s.db.Save(l)
}

func (s shareBackend) Delete(hash string) error {
	err := s.db.DeleteStruct(&share.Link{Hash: hash})
	if err == storm.ErrNotFound {
		return nil
	}
	return err
}

func (s shareBackend) IncAccessCount(hash string) error {
	var v []*share.Link
	err := s.db.Select(q.Eq("Hash", hash)).Find(&v)
	if err == storm.ErrNotFound || len(v) == 0 {
		return nil
	}

	for _, l := range v {
		err = s.db.UpdateField(l, "AccessCount", l.AccessCount+1)
		if err != nil {
			return err
		}
	}
	return s.db.Commit()
}

func (s shareBackend) IncDownloadCount(hash string) error {
	var v []*share.Link
	err := s.db.Select(q.Eq("Hash", hash)).Find(&v)
	if err == storm.ErrNotFound || len(v) == 0 {
		return nil
	}

	for _, l := range v {
		err = s.db.UpdateField(l, "DownloadCount", l.DownloadCount+1)
		if err != nil {
			return err
		}
	}
	return s.db.Commit()
}

