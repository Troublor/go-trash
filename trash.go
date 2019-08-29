package main

import (
	"flag"
	"fmt"
	"github.com/Troublor/trash-go/operation"
	"os"
	"strings"
)

func main() {
	rmCmd := flag.NewFlagSet("operation", flag.ExitOnError)
	rmDirectory := rmCmd.Bool("d", false, "Remove a directory")
	rmRecursive := rmCmd.Bool("r", false, "Recursively remove a directory")

	urCmd := flag.NewFlagSet("ur", flag.ExitOnError)
	urId := urCmd.Bool("i", false, "Use id of the item to retrieve (un-remove) item from trash bin")
	urOverride := urCmd.Bool("o", false, "Override the existing file when retrieve (un-remove) items")

	lsCmd := flag.NewFlagSet("ls", flag.ExitOnError)
	lsDetail := lsCmd.Bool("d", false, "List the detail of all items in trash bin")

	if len(os.Args) < 2 {
		fmt.Println(Usage())
		os.Exit(-1)
	}
	switch os.Args[1] {
	case "rm":
		err := rmCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Usage Error")
			fmt.Println(Usage())
		}
		itemPath := rmCmd.Arg(0)
		id, err := operation.Remove(itemPath, *rmDirectory, *rmRecursive)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("remove " + itemPath + " complete, trash id = " + id)
		}
	case "ur":
		err := urCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Usage Error")
			fmt.Println(Usage())
		}
		payload := urCmd.Arg(0)
		trashInfo, err := operation.UnRemove(payload, *urId, *urOverride)
		if err != nil {
			switch err {
			case operation.ItemNotExistError:
				fmt.Println("can not find " + payload + " in trash bin")
			case operation.ItemExistError:
				fmt.Println("a file or directory already exists in original path, please try again with option -o")
			case operation.MultipleItemsError:
				fmt.Println("multiple items named '" + payload + "' found in trash bin, please specify trash id to retrieve")
			default:
				fmt.Println("retrieve failed")
			}
		} else {
			fmt.Printf("retrieve %s to %s\n", trashInfo.BaseName, trashInfo.OriginalPath)
		}
	case "ls":
		err := lsCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Usage Error")
			fmt.Println(Usage())
		}
		results := operation.List(*lsDetail)
		for _, result := range results {
			fmt.Println(strings.Join(result, "\t"))
		}
	}
}

func Usage() string {
	return "Usage: "
}
