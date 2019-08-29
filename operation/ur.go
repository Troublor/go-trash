package operation

import (
	"errors"
	"os"
)

func UnRemove(payload string, isId bool, override bool) (*TrashInfo, error) {
	if isId {
		return unRemoveById(payload, override)
	} else {
		return unRemoveByName(payload, override)
	}
}

func unRemoveById(id string, override bool) (*TrashInfo, error) {
	trashInfo, err := DbGetTrashItemById(id)
	if err != nil {
		return trashInfo, ItemNotExistError
	}
	// move the item out of trash directory
	if _, err = os.Stat(trashInfo.OriginalPath); err == nil {
		// original path already exist another file
		if override {
			if trashInfo.ItemType == TYPE_DIRECTORY {
				_ = os.RemoveAll(trashInfo.OriginalPath)
			} else if trashInfo.ItemType == TYPE_FILE {
				_ = os.Remove(trashInfo.OriginalPath)
			} else {
				panic(errors.New("invalid argument itemType: " + trashInfo.ItemType))
			}
		} else {
			return trashInfo, ItemExistError
		}
	}
	// delete information in database
	err = DbDeleteTrashItem(id)
	if err != nil {
		return trashInfo, ItemNotExistError
	}
	err = SafeRename(trashInfo.TrashPath, trashInfo.OriginalPath)
	if err != nil {
		panic(err)
	}
	return trashInfo, nil
}
func unRemoveByName(name string, override bool) (*TrashInfo, error) {
	count := 0
	var id string
	items := DbListAllTrashItems()
	for _, item := range items {
		if item.BaseName == name {
			count += 1
			id = item.Id
		}
	}
	if count == 0 {
		return nil, ItemNotExistError
	} else if count > 1 {
		return nil, MultipleItemsError
	} else {
		return unRemoveById(id, override)
	}
}
