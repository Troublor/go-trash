package cmd

import (
	"fmt"
	"github.com/Troublor/go-trash/storage/model"
	"github.com/Troublor/go-trash/system"
	"github.com/spf13/cobra"
	"path"
	"sort"
	"strings"
)

var verboseLs bool

var lsCmd = &cobra.Command{
	Use:   "ls [-v]",
	Short: "List the trashes in the trash bin",
	Long: `List the trashes in the trash bin, users can only view the trashes owned by themselves 
			root user can view all trashes of all users. `,
	Run: func(cmd *cobra.Command, args []string) {
		results := List()
		fmt.Println(results.String(verboseLs))
	},
}

func init() {
	lsCmd.Flags().BoolVarP(&verboseLs, "verbose", "v", false,
		"List the detailed information of trash")
}

func List() model.TrashMetadataList {
	results := db.ListTrashItems(system.GetUser())
	sort.Slice(results, func(i, j int) bool {
		base1, base2 := path.Base(results[i].OriginalPath), path.Base(results[j].OriginalPath)
		return strings.Compare(base1, base2) < 0
	})
	return results
}
