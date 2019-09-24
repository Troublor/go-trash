package operation

import (
	"github.com/Troublor/trash-go/storage"
	"regexp"
	"strings"
)

func Search(keyword string) storage.TrashInfoList {
	re, err := regexp.Compile(keyword)
	if err != nil {
		return searchWithPlainString(keyword)
	}
	return searchWithRegexp(re)
}

func searchWithPlainString(keyword string) storage.TrashInfoList {
	results := storage.DbListAllTrashItems()
	sResults := make([]storage.TrashInfo, 0)
	for _, elem := range results {
		if strings.Index(elem.BaseName, keyword) >= 0 {
			sResults = append(sResults, elem)
		}
	}
	return sResults
}

func searchWithRegexp(re *regexp.Regexp) storage.TrashInfoList {
	results := storage.DbListAllTrashItems()
	sResults := make([]storage.TrashInfo, 0)
	for _, elem := range results {
		if re.MatchString(elem.BaseName) {
			sResults = append(sResults, elem)
		}
	}
	return sResults
}
