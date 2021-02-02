package http

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/favorite"
)

var favoriteListHandler = func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	var (
		s   []*favorite.Favorite
		err error
	)
	s, err = d.store.Favorite.FindByUserID(d.user.ID)
	if err == errors.ErrNotExist {
		return renderJSON(w, r, []*share.Link{})
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return renderJSON(w, r, s)
}

var favoriteDeleteHandler = func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	hash := strings.TrimSuffix(r.URL.Path, "/")
	hash = strings.TrimPrefix(hash, "/")

	if hash == "" {
		return http.StatusBadRequest, nil
	}

	err := d.store.Favorite.Delete(hash, d.user.ID)
	return errToStatus(err), err
}

var favoritePostHandler = func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	var s *favorite.Favorite

	if r.Body != nil {
		defer r.Body.Close()
	}

	bytes := make([]byte, 6)
	_, err := rand.Read(bytes)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	hash := base64.URLEncoding.EncodeToString(bytes)

	s = &favorite.Favorite{
		Path:   r.URL.Path,
		Hash:   hash,
		UserID: d.user.ID,
		Name:   "", // TODO: get name from path
	}

	if err := d.store.Favorite.Save(s); err != nil {
		return http.StatusInternalServerError, err
	}

	return renderJSON(w, r, s)
}
