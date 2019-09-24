package storage

import (
	"github.com/bndr/gotabulate"
	"time"
)

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

type TrashInfoList []TrashInfo

func (list TrashInfoList) ToString(detailed bool) string {
	if len(list) == 0 {
		return "No data found"
	}
	var payload [][]interface{}
	for _, elem := range list {
		if detailed {
			payload = append(payload, []interface{}{
				elem.Id,
				elem.BaseName,
				elem.OriginalPath,
				elem.ItemType,
				elem.DeleteTime,
				elem.TrashPath,
			})
		} else {
			payload = append(payload, []interface{}{
				elem.Id,
				elem.BaseName,
				elem.OriginalPath,
			})
		}
	}
	t := gotabulate.Create(payload)
	if detailed {
		t.SetHeaders([]string{"Index", "Basename", "Original Path", "Type", "Delete Time", "Trash Path"})
	} else {
		t.SetHeaders([]string{"Index", "Basename", "Original Path"})
	}
	t.SetAlign("left")
	return t.Render("simple")
}
