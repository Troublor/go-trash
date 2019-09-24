package operation

import (
	"github.com/Troublor/trash-go/storage"
	"path"
	"sort"
	"strings"
)

func List(containDetail bool) [][]string {
	results := storage.DbListAllTrashItems()
	sort.Slice(results, func(i, j int) bool {
		base1, base2 := path.Base(results[i].OriginalPath), path.Base(results[j].OriginalPath)
		return strings.Compare(base1, base2) < 0
	})
	r := make([][]string, len(results))
	for i, result := range results {
		if containDetail {
			r[i] = []string{result.Id, result.OriginalPath, result.TrashPath, result.BaseName, result.ItemType, result.DeleteTime.Format("2006-01-02 15:04:05")}
		} else {
			r[i] = []string{result.Id, result.BaseName, result.OriginalPath}
		}
	}
	return r
}
