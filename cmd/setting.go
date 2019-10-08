package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var settingCmd = &cobra.Command{
	Use:   "setting",
	Short: "Sub-command setting, which does operations on settings of Go-trash",
	Long:  `Sub-command setting, which does operations on settings of Go-trash`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	},
}

func init() {
	settingCmd.AddCommand(settingLsCmd)
}
