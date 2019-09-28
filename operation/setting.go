package operation

import "github.com/Troublor/trash-go/storage"

func ListSettings() storage.SettingList {
	return storage.ListAllSettings()
}
