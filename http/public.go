package http

import (
	"errors"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"

	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/share"
)

var withHashFile = func(fn handleFunc) handleFunc {
	return func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		id, path := ifPathWithName(r)
		link, err := d.store.Share.GetByHash(id)
		if err != nil {
			return errToStatus(err), err
		}

		status, err := authenticateShareRequest(r, link, d)
		if err != nil {
			return status, err
		}

		if path == "/" {
			d.store.Share.IncAccectCount(link.Hash)
		} else {
			d.store.Share.IncDownloadCount(link.Hash)
		}

		user, err := d.store.Users.Get(d.server.Root, link.UserID)
		if err != nil {
			return errToStatus(err), err
		}

		d.user = user

		file, err := files.NewFileInfo(files.FileOptions{
			Fs:              d.user.Fs,
			Path:            link.Path,
			Modify:          d.user.Perm.Modify,
			Expand:          true,
			ReadHeader:      d.server.TypeDetectionByHeader,
			Checker:         d,
			SharedCodeToken: link.SharedCodeToken,
		})
		if err != nil {
			return errToStatus(err), err
		}

		if file.IsDir {
			// set fs root to the shared folder
			d.user.Fs = afero.NewBasePathFs(d.user.Fs, filepath.Dir(link.Path))

			file, err = files.NewFileInfo(files.FileOptions{
				Fs:              d.user.Fs,
				Path:            path,
				Modify:          d.user.Perm.Modify,
				Expand:          true,
				Checker:         d,
				SharedCodeToken: link.SharedCodeToken,
			})
			if err != nil {
				return errToStatus(err), err
			}
		}

		d.raw = file
		return fn(w, r, d)
	}
}

// ref to https://github.com/filebrowser/filebrowser/pull/727
// `/api/public/dl/MEEuZK-v/file-name.txt` for old browsers to save file with correct name
func ifPathWithName(r *http.Request) (id, filePath string) {
	pathElements := strings.Split(r.URL.Path, "/")
	// prevent maliciously constructed parameters like `/api/public/dl/XZzCDnK2_not_exists_hash_name`
	// len(pathElements) will be 1, and golang will panic `runtime error: index out of range`

	switch len(pathElements) {
	case 1:
		return r.URL.Path, "/"
	default:
		return pathElements[0], path.Join("/", path.Join(pathElements[1:]...))
	}
}

var publicShareHandler = withHashFile(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	file := d.raw.(*files.FileInfo)

	if file.IsDir {
		file.Listing.Sorting = files.Sorting{By: "name", Asc: false}
		file.Listing.ApplySort()
		return renderJSON(w, r, file)
	}

	return renderJSON(w, r, file)
})

var publicDlHandler = withHashFile(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	file := d.raw.(*files.FileInfo)
	if !file.IsDir {
		return rawFileHandler(w, r, file)
	}

	return rawDirHandler(w, r, d, file)
})

func authenticateShareRequest(r *http.Request, l *share.Link, d *data) (int, error) {
	if l.SharedCode == "" {
		return 0, nil
	}

	sharedCodeToken := r.URL.Query().Get("shared_code_token")
	if sharedCodeToken == l.SharedCodeToken {
		return 0, nil
	}

	shared_code := r.Header.Get("X-SHARED-CODE")
	if l.SharedCode != shared_code {
		return http.StatusUnauthorized, errors.New("Invalid shared code")
	}
	return 0, nil
}
