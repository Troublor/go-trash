package model

import (
	"github.com/bndr/gotabulate"
	"time"
)

const (
	TYPE_DIRECTORY = "D"
	TYPE_FILE      = "F"
)

type TrashMetadata struct {
	ID           string `gorm:"primary_key"`
	OriginalPath string // the original path of the item
	TrashPath    string // the path in the trash bin
	BaseName     string // the base name of the item
	Type         string // file type of the item
	Owner        string // owner user of the item
	CreatedAt    time.Time
}

func (trashInfo TrashMetadata) IsDirectory() bool {
	return trashInfo.Type == TYPE_DIRECTORY
}

func (trashInfo TrashMetadata) IsFile() bool {
	return trashInfo.Type == TYPE_FILE
}

func (trashInfo TrashMetadata) OwnedBy(owner string) bool {
	return trashInfo.Owner == owner
}

func (trashInfo TrashMetadata) Equals(info TrashMetadata) bool {
	return trashInfo.OriginalPath == info.OriginalPath &&
		trashInfo.CreatedAt.Equal(info.CreatedAt)
}

type TrashMetadataList []TrashMetadata

func (list TrashMetadataList) Contains(info TrashMetadata) bool {
	for _, trash := range list {
		if trash.Equals(info) {
			return true
		}
	}
	return false
}

func (list *TrashMetadataList) Merge(anotherList TrashMetadataList) {
	for _, info := range anotherList {
		if !list.Contains(info) {
			*list = append(*list, info)
		}
	}
}

func (list TrashMetadataList) String(detailed bool) string {
	if len(list) == 0 {
		return "No data found"
	}
	var payload [][]interface{}
	for _, elem := range list {
		if detailed {
			payload = append(payload, []interface{}{
				elem.ID,
				elem.BaseName,
				elem.OriginalPath,
				elem.Type,
				elem.Owner,
				elem.CreatedAt,
				elem.TrashPath,
			})
		} else {
			payload = append(payload, []interface{}{
				elem.ID,
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
