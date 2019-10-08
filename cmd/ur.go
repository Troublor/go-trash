package cmd

import (
	"errors"
	"fmt"
	"github.com/Troublor/trash-go/errs"
	"github.com/Troublor/trash-go/storage"
	"github.com/Troublor/trash-go/system"
	"github.com/spf13/cobra"
	"os"
)

var id bool
var override bool

var urCmd = &cobra.Command{
	Use:   "ur [-i]|[-o]",
	Short: "Un-remove: retrieve files or directories from trash bin",
	Long:  `Un-remove: retrieve files or directories from trash bin`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, payload := range args {
			trashInfo, err := UnRemove(payload, id, override)
			if err != nil {
				switch err {
				case errs.ItemNotExistError:
					fmt.Println("UnRemove Error: " + "can not find " + payload + " in trash bin")
				case errs.ItemExistError:
					fmt.Println("UnRemove Error: " + "a file or directory already exists in original path of " + payload + ", please try again with option -o")
				case errs.MultipleItemsError:
					fmt.Println("UnRemove Error: " + "multiple items named '" + payload + "' found in trash bin, please specify trash id to retrieve")
				default:
					fmt.Println("UnRemove Error: " + "retrieve failed")
				}
			} else {
				fmt.Printf("retrieve %s to %s\n", trashInfo.BaseName, trashInfo.OriginalPath)
			}
		}
	},
}

func init() {
	urCmd.Flags().BoolVarP(&id, "id", "i", false,
		"use id of the item to retrieve (un-remove) item from trash bin")
	urCmd.Flags().BoolVarP(&override, "override", "o", false,
		"override the existing file when retrieve (un-remove) items")
}

func UnRemove(payload string, isId bool, override bool) (*storage.TrashInfo, error) {
	if isId {
		return unRemoveById(payload, override)
	} else {
		return unRemoveByName(payload, override)
	}
}

func unRemoveById(id string, override bool) (*storage.TrashInfo, error) {
	trashInfo, err := storage.DbGetTrashItemById(id, system.GetUser())
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
	err = storage.DbDeleteTrashItem(id, system.GetUser())
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
	items := storage.DbListAllTrashItems(system.GetUser())
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
