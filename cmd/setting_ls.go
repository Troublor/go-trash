package cmd

import (
	"fmt"
	"github.com/Troublor/trash-go/storage"
	"github.com/spf13/cobra"
)

var settingLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Sub-command setting, which does operations on settings of Go-trash",
	Long:  `Sub-command setting, which does operations on settings of Go-trash`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
	},
}

func ListSettings() storage.SettingList {
	return storage.ListAllSettings()
}
