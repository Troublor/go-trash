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
	Owner        string
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

func (trashInfo TrashInfo) OwnedBy(owner string) bool {
	return trashInfo.Owner == owner
}

func (trashInfo TrashInfo) Equals(info TrashInfo) bool {
	return trashInfo.OriginalPath == info.OriginalPath &&
		trashInfo.DeleteTime.Equal(info.DeleteTime)
}

type TrashInfoList []TrashInfo

func (list TrashInfoList) Contains(info TrashInfo) bool {
	for _, trash := range list {
		if trash.Equals(info) {
			return true
		}
	}
	return false
}

func (list *TrashInfoList) Merge(anotherList TrashInfoList) {
	for _, info := range anotherList {
		if !list.Contains(info) {
			*list = append(*list, info)
		}
	}
}

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
				elem.Owner,
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
		t.SetHeaders([]string{"Index", "Basename", "Original Path", "Type", "Owner", "Delete Time", "Trash Path"})
	} else {
		t.SetHeaders([]string{"Index", "Basename", "Original Path"})
	}
	t.SetAlign("left")
	return t.Render("simple")
}
