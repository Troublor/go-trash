package cmd

import (
	"fmt"
	"github.com/Troublor/go-trash/storage/model"
	"github.com/Troublor/go-trash/system"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
)

var verboseSs bool

var ssCmd = &cobra.Command{
	Use:   "ss [-v]",
	Short: "Search trash in trash bin",
	Long:  `Search trash in trash bin`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		results := model.TrashMetadataList{}
		for _, kw := range args {
			results.Merge(Search(kw))
		}
		fmt.Println(results.String(verboseSs))
	},
}

func init() {
	ssCmd.Flags().BoolVarP(&verboseSs, "verbose", "v", false,
		"show the detail of searched items in trash bin")
}

func Search(keyword string) model.TrashMetadataList {
	re, err := regexp.Compile(keyword)
	if err != nil {
		return searchWithPlainString(keyword)
	}
	return searchWithRegexp(re)
}

func searchWithPlainString(keyword string) model.TrashMetadataList {
	results := db.ListTrashItems(system.GetUser())
	sResults := make([]model.TrashMetadata, 0)
	for _, elem := range results {
		if strings.Index(elem.BaseName, keyword) >= 0 {
			sResults = append(sResults, elem)
		}
	}
	return sResults
}

func searchWithRegexp(re *regexp.Regexp) model.TrashMetadataList {
	results := db.ListTrashItems(system.GetUser())
	sResults := make([]model.TrashMetadata, 0)
	for _, elem := range results {
		if re.MatchString(elem.BaseName) {
			sResults = append(sResults, elem)
		}
	}
	return sResults
}
