package cmd

import (
	"fmt"
	"github.com/Troublor/trash-go/errs"
	"github.com/Troublor/trash-go/storage"
	"github.com/Troublor/trash-go/system"
	"github.com/spf13/cobra"
	"os"
)

var useId bool
var clearAll bool

var cleanCmd = &cobra.Command{
	Use:   "clean [-i][-a] items...",
	Short: "Clean the trash bin, permanently delete items in trash bin.",
	Long:  `Clean the items listed in the command, if option -i is used, items should be given by their indices.`,
	Run: func(cmd *cobra.Command, args []string) {
		if clearAll {
			// delete all items in trash bin
			items := List()
			for _, item := range items {
				err := PermanentlyDelete(item.Id)
				if err != nil {
					fmt.Printf("Clean Error: %s", err.Error())
				}
			}
			return
		}
		for _, arg := range args {
			err := Clean(useId, arg)
			if err != nil {
				switch err {
				case errs.ItemNotExistError:
					fmt.Printf("Clean Error: item \"%s\" does not exist\n", arg)
				case errs.MultipleItemsError:
					fmt.Printf("Clean Error: multiple items named \"%s\" found in trash bin, please specify trash id to retrieve", arg)
				default:
					fmt.Printf("Clean Error: %s", err.Error())
				}
			}
		}
	},
}

func init() {
	cleanCmd.Flags().BoolVarP(&useId, "id", "i", false,
		"use id of the item to clean (permanently delete) item from trash bin")
	cleanCmd.Flags().BoolVarP(&clearAll, "all", "a", false,
		"clean all the trash items, i.e. permanently delete all the items in trash bin.")
}

func Clean(useId bool, item string) error {
	results := List()
	if useId {
		return PermanentlyDelete(item)
	} else {
		count := 0
		var id string
		for _, result := range results {
			if item == result.BaseName {
				count++
				id = result.Id
			}
		}
		if count > 1 {
			return errs.MultipleItemsError
		} else if count == 0 {
			return errs.ItemNotExistError
		} else {
			return PermanentlyDelete(id)
		}
	}
}

func PermanentlyDelete(id string) error {
	results := List()
	for _, item := range results {
		if item.Id == id {
			err := os.RemoveAll(item.TrashPath)
			if err == nil {
				err = storage.DbDeleteTrashItem(id, system.GetUser())
			}
			return err
		}
	}
	return errs.ItemNotExistError
}
