package trash

type Trash struct {
	Hash       string `json:"hash" storm:"id,index"`
	OriginPath string `json:"originPath" storm:"index"`
	TrashPath  string `json:"trashPath" storm:"index"`
	UserID     uint   `json:"userID"`
	Datetime   int64  `json:"datetime"`
}
