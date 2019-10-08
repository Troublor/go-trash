package cmd

import (
	"fmt"
	"github.com/Troublor/trash-go/storage"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Go-trash",
	Long:  `Print the version number of Go-trash`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(storage.VersionString())
	},
}
