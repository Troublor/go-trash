package operation

import (
	"errors"
	"github.com/Troublor/trash-go/errs"
	"github.com/Troublor/trash-go/storage"
	"os"
)

func UnRemove(payload string, isId bool, override bool) (*storage.TrashInfo, error) {
	if isId {
		return unRemoveById(payload, override)
	} else {
		return unRemoveByName(payload, override)
	}
}

func unRemoveById(id string, override bool) (*storage.TrashInfo, error) {
	trashInfo, err := storage.DbGetTrashItemById(id)
	if err != nil {
		return trashInfo, errs.ItemNotExistError
	}
	// move the item out of trash directory
	if _, err = os.Stat(trashInfo.OriginalPath); err == nil {
		// original path already exist another file
		if override {
			if trashInfo.ItemType == storage.TYPE_DIRECTORY {
				_ = os.RemoveAll(trashInfo.OriginalPath)
			} else if trashInfo.ItemType == storage.TYPE_FILE {
				_ = os.Remove(trashInfo.OriginalPath)
			} else {
				panic(errors.New("invalid argument itemType: " + trashInfo.ItemType))
			}
		} else {
			return trashInfo, errs.ItemExistError
		}
	}
	// delete information in database
	err = storage.DbDeleteTrashItem(id)
	if err != nil {
		return trashInfo, errs.ItemNotExistError
	}
	err = storage.SafeRename(trashInfo.TrashPath, trashInfo.OriginalPath)
	if err != nil {
		panic(err)
	}
	return trashInfo, nil
}
func unRemoveByName(name string, override bool) (*storage.TrashInfo, error) {
	count := 0
	var id string
	items := storage.DbListAllTrashItems()
	for _, item := range items {
		if item.BaseName == name {
			count += 1
			id = item.Id
		}
	}
	if count == 0 {
		return nil, errs.ItemNotExistError
	} else if count > 1 {
		return nil, errs.MultipleItemsError
	} else {
		return unRemoveById(id, override)
	}
}
