package cmd

import (
	"fmt"
	"github.com/Troublor/trash-go/storage"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "gotrash",
	Short: "Go-trash is a linux command-line trash files management tool",
	Long: `Go-trash is a linux command-line trash files management tool which provides
			 features similar to the Recycle Bin in Windows.
		   Developed by Troublor, 2019`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initialize)

	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(settingCmd)
	rootCmd.AddCommand(ssCmd)
	rootCmd.AddCommand(urCmd)
	rootCmd.AddCommand(versionCmd)
}

func initialize() {
	storage.InitStorage()
}
