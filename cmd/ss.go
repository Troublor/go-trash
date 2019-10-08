package cmd

import (
	"fmt"
	"github.com/Troublor/trash-go/storage"
	"github.com/Troublor/trash-go/system"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
)

var verboseSs bool

var ssCmd = &cobra.Command{
	Use:   "ss [-v]",
	Short: "Search trash in trash bin",
	Long:  `Search trash in trash bin`,
	Run: func(cmd *cobra.Command, args []string) {
		results := &storage.TrashInfoList{}
		for _, kw := range args {
			results.Merge(Search(kw))
		}
		fmt.Println(results.ToString(verboseSs))
	},
}

func init() {
	ssCmd.Flags().BoolVarP(&verboseSs, "verbose", "v", false,
		"show the detail of searched items in trash bin")
}

func Search(keyword string) storage.TrashInfoList {
	re, err := regexp.Compile(keyword)
	if err != nil {
		return searchWithPlainString(keyword)
	}
	return searchWithRegexp(re)
}

func searchWithPlainString(keyword string) storage.TrashInfoList {
	results := storage.DbListAllTrashItems(system.GetUser())
	sResults := make([]storage.TrashInfo, 0)
	for _, elem := range results {
		if strings.Index(elem.BaseName, keyword) >= 0 {
			sResults = append(sResults, elem)
		}
	}
	return sResults
}

func searchWithRegexp(re *regexp.Regexp) storage.TrashInfoList {
	results := storage.DbListAllTrashItems(system.GetUser())
	sResults := make([]storage.TrashInfo, 0)
	for _, elem := range results {
		if re.MatchString(elem.BaseName) {
			sResults = append(sResults, elem)
		}
	}
	return sResults
}
