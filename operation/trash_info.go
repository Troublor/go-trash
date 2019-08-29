package operation

import "time"

type TrashInfo struct {
	Id           string
	OriginalPath string
	TrashPath    string
	BaseName     string
	ItemType     string
	DeleteTime   time.Time
}

const (
	TYPE_FILE      string = "F"
	TYPE_DIRECTORY string = "D"
)

func (trashInfo TrashInfo) isDirectory() bool {
	return trashInfo.ItemType == TYPE_DIRECTORY
}

func (trashInfo TrashInfo) isFile() bool {
	return trashInfo.ItemType == TYPE_FILE
}
