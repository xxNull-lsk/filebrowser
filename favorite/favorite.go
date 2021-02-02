package favorite

// Favorite is the information needed to build a favorite item.
type Favorite struct {
	Hash   string `json:"hash" storm:"id,index"`
	Path   string `json:"path" storm:"index"`
	UserID uint   `json:"userID"`
	Name   uint   `json:"name"`
	Type   uint   `json:"type"`
}
