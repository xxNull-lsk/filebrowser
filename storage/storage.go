package storage

import (
	"github.com/filebrowser/filebrowser/v2/auth"
	"github.com/filebrowser/filebrowser/v2/favorite"
	"github.com/filebrowser/filebrowser/v2/settings"
	"github.com/filebrowser/filebrowser/v2/share"
	"github.com/filebrowser/filebrowser/v2/trash"
	"github.com/filebrowser/filebrowser/v2/users"
)

// Storage is a storage powered by a Backend which makes the necessary
// verifications when fetching and saving data to ensure consistency.
type Storage struct {
	Users    *users.Storage
	Share    *share.Storage
	Trash    *trash.Storage
	Favorite *favorite.Storage
	Auth     *auth.Storage
	Settings *settings.Storage
}
