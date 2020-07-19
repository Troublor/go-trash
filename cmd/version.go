package cmd

import (
	"fmt"
	"github.com/Troublor/go-trash/storage"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Go-trash",
	Long:  `Print the version number of Go-trash`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(storage.Version())
	},
}
