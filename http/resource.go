package http

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/afero"

	"github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/fileutils"
	"github.com/filebrowser/filebrowser/v2/rules"
	"github.com/filebrowser/filebrowser/v2/trash"
)

func IsTrash(r *http.Request) bool {
	return strings.TrimRight(r.URL.Path, "/") == "/.trash"
}

var resourceGetHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	isTrash := IsTrash(r)
	if isTrash && !files.Exist(d.user.Fs, r.URL.Path) {
		err := d.user.Fs.Mkdir("/.trash", 0777)
		if err != nil {
			return errToStatus(err), err
		}
	}

	file, err := files.NewFileInfo(files.FileOptions{
		Fs:         d.user.Fs,
		Path:       r.URL.Path,
		Modify:     d.user.Perm.Modify,
		Expand:     true,
		ReadHeader: d.server.TypeDetectionByHeader,
		Checker:    d,
	})
	if err != nil {
		fmt.Println(err)
		return errToStatus(err), err
	}

	if file.IsDir {
		if d.user.HideDotfiles && !isTrash {
			count := len(file.Listing.Items)
			for i := range file.Listing.Items {
				index := count - i - 1
				item := file.Listing.Items[index]
				if rules.MatchHidden(item.Path) {
					file.Listing.Items = append(file.Listing.Items[:index], file.Listing.Items[index+1:]...)
					if item.IsDir {
						file.Listing.NumDirs--
					} else {
						file.Listing.NumFiles--
					}
				}
			}
		}
		if isTrash {
			count := len(file.Listing.Items)
			for i := range file.Listing.Items {
				index := count - i - 1
				item := file.Listing.Items[index]
				t, err := d.store.Trash.GetByHash(item.Name)
				if err != nil || t.UserID != d.user.ID && !d.user.Perm.Admin {
					file.Listing.Items = append(file.Listing.Items[:index], file.Listing.Items[index+1:]...)
					continue
				}

				item.Name = filepath.Base(t.OriginPath)
				item.OriginPath = t.OriginPath
				item.DeleteTime = time.Unix(t.Datetime, 0)
				user, err := d.store.Users.Get("", t.UserID)
				if err == nil {
					item.Username = user.Username
				}
			}
		}
		file.Listing.Sorting = d.user.Sorting
		file.Listing.ApplySort()
		return renderJSON(w, r, file)
	}

	if checksum := r.URL.Query().Get("checksum"); checksum != "" {
		err := file.Checksum(checksum)
		if err == errors.ErrInvalidOption {
			return http.StatusBadRequest, nil
		} else if err != nil {
			return http.StatusInternalServerError, err
		}

		// do not waste bandwidth if we just want the checksum
		file.Content = ""
	}

	return renderJSON(w, r, file)
})

func resourceTrashHandler(file *files.FileInfo, w http.ResponseWriter, r *http.Request, d *data) (int, error) {

	bytes := make([]byte, 6)
	_, err := rand.Read(bytes)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	hash := base64.URLEncoding.EncodeToString(bytes)

	dst := "/.trash/" + hash
	if !files.Exist(d.user.Fs, "/.trash") {
		err = d.user.Fs.Mkdir("/.trash", 0777)
		if err != nil {
			return errToStatus(err), err
		}
	}

	s := &trash.Trash{
		Hash:       hash,
		OriginPath: file.Path,
		TrashPath:  dst,
		UserID:     d.user.ID,
		Datetime:   time.Now().Unix(),
	}
	err = d.store.Trash.Save(s)
	if err != nil {
		return errToStatus(err), err
	}

	err = fileutils.MoveFile(d.user.Fs, file.Path, dst)
	if err != nil {
		return errToStatus(err), err
	}

	return http.StatusOK, nil
}

var trashDeleteHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	if r.URL.Path == "/" || !d.user.Perm.Delete {
		return http.StatusForbidden, nil
	}
	pathElements := strings.Split(r.URL.Path, "/")
	if len(pathElements) == 0 {
		return http.StatusForbidden, nil
	}
	hash := pathElements[len(pathElements)-1]
	if hash == "" {
		return http.StatusBadRequest, os.ErrInvalid
	}

	s, err := d.store.Trash.GetByHash(hash)
	if err != nil {
		return errToStatus(err), err
	}

	_, err = d.user.Fs.Stat(s.OriginPath)
	if err == nil {
		// if s.OriginPath is exist file/folder, can't restore it.
		return http.StatusFound, err
	}
	err = fileutils.MoveFile(d.user.Fs, s.TrashPath, s.OriginPath)
	if err != nil {
		return errToStatus(err), err
	}

	err = d.store.Trash.Delete(hash)
	if err != nil {
		return errToStatus(err), err
	}

	return http.StatusOK, nil
})

func resourceDeleteHandler(fileCache FileCache) handleFunc {
	return withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		if r.URL.Path == "/" || !d.user.Perm.Delete {
			return http.StatusForbidden, nil
		}

		file, err := files.NewFileInfo(files.FileOptions{
			Fs:         d.user.Fs,
			Path:       r.URL.Path,
			Modify:     d.user.Perm.Modify,
			Expand:     true,
			ReadHeader: d.server.TypeDetectionByHeader,
			Checker:    d,
		})
		if err != nil {
			return errToStatus(err), err
		}

		action := r.URL.Query().Get("action")
		if action == "trash" {
			return resourceTrashHandler(file, w, r, d)
		}

		// delete thumbnails
		for _, previewSizeName := range PreviewSizeNames() {
			size, _ := ParsePreviewSize(previewSizeName)
			if err := fileCache.Delete(r.Context(), previewCacheKey(file.Path, size)); err != nil { //nolint:govet
				return errToStatus(err), err
			}
		}

		err = d.RunHook(func() error {
			return d.user.Fs.RemoveAll(r.URL.Path)
		}, "delete", r.URL.Path, "", d.user)

		if err != nil {
			return errToStatus(err), err
		}

		return http.StatusOK, nil
	})
}

var resourcePostPutHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	if !d.user.Perm.Create && r.Method == http.MethodPost {
		return http.StatusForbidden, nil
	}

	if !d.user.Perm.Modify && r.Method == http.MethodPut {
		return http.StatusForbidden, nil
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, r.Body)
	}()

	// For directories, only allow POST for creation.
	if strings.HasSuffix(r.URL.Path, "/") {
		if r.Method == http.MethodPut {
			return http.StatusMethodNotAllowed, nil
		}
		err := d.user.Fs.MkdirAll(r.URL.Path, 0777)
		return errToStatus(err), err
	}

	if r.Method == http.MethodPost && r.URL.Query().Get("override") != "true" {
		if _, err := d.user.Fs.Stat(r.URL.Path); err == nil {
			return http.StatusConflict, nil
		}
	}

	action := "upload"
	if r.Method == http.MethodPut {
		action = "save"
	}

	err := d.RunHook(func() error {
		dir, _ := path.Split(r.URL.Path)
		err := d.user.Fs.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}

		file, err := d.user.Fs.OpenFile(r.URL.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, r.Body)
		if err != nil {
			return err
		}

		// Gets the info about the file.
		info, err := file.Stat()
		if err != nil {
			return err
		}

		etag := fmt.Sprintf(`"%x%x"`, info.ModTime().UnixNano(), info.Size())
		w.Header().Set("ETag", etag)
		return nil
	}, action, r.URL.Path, "", d.user)

	if err != nil {
		_ = d.user.Fs.RemoveAll(r.URL.Path)
	}

	return errToStatus(err), err
})

var resourcePatchHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	src := r.URL.Path
	dst := r.URL.Query().Get("destination")
	action := r.URL.Query().Get("action")
	dst, err := url.QueryUnescape(dst)
	if err != nil {
		return errToStatus(err), err
	}
	if dst == "/" || src == "/" {
		return http.StatusForbidden, nil
	}
	if err = checkParent(src, dst); err != nil {
		return http.StatusBadRequest, err
	}

	override := r.URL.Query().Get("override") == "true"
	rename := r.URL.Query().Get("rename") == "true"
	if !override && !rename {
		if _, err = d.user.Fs.Stat(dst); err == nil {
			return http.StatusConflict, nil
		}
	}
	if rename {
		dst = addVersionSuffix(dst, d.user.Fs)
	}

	err = d.RunHook(func() error {
		switch action {
		// TODO: use enum
		case "copy":
			if !d.user.Perm.Create {
				return errors.ErrPermissionDenied
			}

			return fileutils.Copy(d.user.Fs, src, dst)
		case "rename":
			if !d.user.Perm.Rename {
				return errors.ErrPermissionDenied
			}
			src = path.Clean("/" + src)
			dst = path.Clean("/" + dst)

			return fileutils.MoveFile(d.user.Fs, src, dst)
		default:
			return fmt.Errorf("unsupported action %s: %w", action, errors.ErrInvalidRequestParams)
		}
	}, action, src, dst, d.user)

	return errToStatus(err), err
})

func checkParent(src, dst string) error {
	rel, err := filepath.Rel(src, dst)
	if err != nil {
		return err
	}

	rel = filepath.ToSlash(rel)
	if !strings.HasPrefix(rel, "../") && rel != ".." && rel != "." {
		return errors.ErrSourceIsParent
	}

	return nil
}

func addVersionSuffix(source string, fs afero.Fs) string {
	counter := 1
	dir, name := path.Split(source)
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)

	for {
		if _, err := fs.Stat(source); err != nil {
			break
		}
		renamed := fmt.Sprintf("%s(%d)%s", base, counter, ext)
		source = path.Join(dir, renamed)
		counter++
	}

	return source
}
