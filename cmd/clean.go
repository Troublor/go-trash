package cmd

import (
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean [-i] items...",
	Short: "Clean the trash bin, permanently delete items in trash bin.",
	Long:  `Clean the items listed in the command, if option -i is used, items should be given by their indices.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
}

func Clean(index bool, items ...string) {
	results := List()

}
