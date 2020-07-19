package cmd

import (
	"errors"
	"fmt"
	"github.com/Troublor/go-trash/errs"
	"github.com/Troublor/go-trash/storage/model"
	"github.com/Troublor/go-trash/system"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var id bool
var override bool
var target string
var parent bool

var urCmd = &cobra.Command{
	Use:   "ur [-i] [-o] [-t target_path [-p]]",
	Short: "Un-remove: retrieve files or directories from trash bin",
	Long: `Un-remove: retrieve files or directories from trash bin, 
If target path is specified, the retrieved files or directories will be put in to the target path. 
If it is not specified, they will be put into original path`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, payload := range args {
			trashInfo, err := UnRemove(payload, id, override, target, parent)
			if err != nil {
				switch err {
				case errs.ItemNotExistError:
					fmt.Println("UnRemove Error: " + "can not find " + payload + " in trash bin")
				case errs.ItemExistError:
					fmt.Println("UnRemove Error: " + "a file or directory already exists in original path of " + payload + ", please try again with option -o")
				case errs.MultipleItemsError:
					fmt.Println("UnRemove Error: " + "multiple items named '" + payload + "' found in trash bin, please specify trash id to retrieve")
				default:
					fmt.Println("UnRemove Error: " + err.Error())
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
	urCmd.Flags().StringVarP(&target, "target", "t", "/original",
		"the directory that the retrieved files or directories will be put into. If not specified, they will be put into their own original path")
	urCmd.Flags().BoolVarP(&parent, "parent", "p", false,
		"no error if the target path does not exist, make parent directories as needed")
}

func AllTrashNames() (names []string) {
	names = make([]string, 0)
	for _, info := range List() {
		var tmp = make([]string, len(names))
		copy(tmp, names)
		for i, name := range tmp {
			if info.BaseName == name {
				names[len(names)-1], names[i] = names[i], names[len(names)-1]
				names = names[:len(names)-1]
			} else {
				names = append(names, info.BaseName)
			}
		}
	}
	return
}

func UnRemove(payload string, isId bool, override bool, target string, parent bool) (model.TrashMetadata, error) {
	if isId {
		return unRemoveById(payload, override, target, parent)
	} else {
		return unRemoveByName(payload, override, target, parent)
	}
}

func unRemoveById(id string, override bool, target string, parent bool) (model.TrashMetadata, error) {
	trashInfo, err := db.GetTrashItemById(id, system.GetUser())
	if err != nil {
		return trashInfo, errs.ItemNotExistError
	}
	// the target path to put the retrieved file or directory
	var targetPath string
	if target == "/original" {
		targetPath = trashInfo.OriginalPath
	} else {
		if _, err = os.Stat(target); err != nil {
			if os.IsNotExist(err) {
				if parent {
					if err = os.MkdirAll(target, os.ModePerm); err != nil {
						return trashInfo, err
					}
				} else {
					return trashInfo, errs.NewFileOrDirNotExistError(target)
				}
			} else {
				return trashInfo, err
			}
		}
		targetPath = filepath.Join(target, trashInfo.BaseName)
	}
	// move the item out of trash directory
	if _, err = os.Stat(targetPath); err == nil {
		// original path already exist another file
		if override {
			if trashInfo.Type == model.TYPE_DIRECTORY {
				_ = os.RemoveAll(targetPath)
			} else if trashInfo.Type == model.TYPE_FILE {
				_ = os.Remove(targetPath)
			} else {
				panic(errors.New("invalid argument itemType: " + trashInfo.Type))
			}
		} else {
			return trashInfo, errs.ItemExistError
		}
	}
	// delete information in database
	err = db.DeleteTrashItem(id, system.GetUser())
	if err != nil {
		return trashInfo, errs.ItemNotExistError
	}
	err = system.SafeRename(trashInfo.TrashPath, targetPath)
	if err != nil {
		panic(err)
	}
	return trashInfo, nil
}
func unRemoveByName(name string, override bool, target string, parent bool) (model.TrashMetadata, error) {
	count := 0
	var id string
	items := db.ListTrashItems(system.GetUser())
	for _, item := range items {
		if item.BaseName == name {
			count += 1
			id = item.ID
		}
	}
	if count == 0 {
		return model.TrashMetadata{}, errs.ItemNotExistError
	} else if count > 1 {
		return model.TrashMetadata{}, errs.MultipleItemsError
	} else {
		return unRemoveById(id, override, target, parent)
	}
}
