package operation

import (
	"github.com/Troublor/trash-go/storage"
	"path"
	"sort"
	"strings"
)

func List() storage.TrashInfoList {
	results := storage.DbListAllTrashItems()
	sort.Slice(results, func(i, j int) bool {
		base1, base2 := path.Base(results[i].OriginalPath), path.Base(results[j].OriginalPath)
		return strings.Compare(base1, base2) < 0
	})
	return results
}
