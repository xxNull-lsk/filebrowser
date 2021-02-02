package http

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
		return renderJSON(w, r, []*favorite.Favorite{})
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return renderJSON(w, r, s)
}

var favoriteGetHandler = func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {

	if r.Body != nil {
		r.Body.Close()
	}

	_, path := ifPathWithName(r)
	fmt.Println(path)

	s, err := d.store.Favorite.Get(path, d.user.ID)
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

	err := d.store.Favorite.Delete(hash)
	return errToStatus(err), err
}

type createfavoriteRequest struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

var favoritePostHandler = func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	var s *favorite.Favorite

	if r.Body == nil {
		return http.StatusInternalServerError, errors.ErrEmptyRequest
	}

	req := &createfavoriteRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	bytes := make([]byte, 6)
	_, err = rand.Read(bytes)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	hash := base64.URLEncoding.EncodeToString(bytes)

	fs, err := d.user.Fs.Stat(req.Path)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	var t string
	if fs.IsDir() {
		t = "folder"
	} else {
		t = "file"
	}

	s = &favorite.Favorite{
		Path:   req.Path,
		Hash:   hash,
		UserID: d.user.ID,
		Name:   req.Name,
		Type:   t,
	}

	if err := d.store.Favorite.Save(s); err != nil {
		return http.StatusInternalServerError, err
	}

	return renderJSON(w, r, s)
}
