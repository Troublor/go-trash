package storage

import (
	"strings"
)

var DefaultSettings = map[string]string{"autoremove": "off"}

type Setting struct {
	key   string
	value string
}

func NewSetting(key, value string) *Setting {
	return &Setting{key: key, value: value}
}

func (setting *Setting) SetValue(value string) error {
	err := UpdateSetting(setting.key, value)
	if err == nil {
		setting.value = value
	}
	return err
}

func (setting Setting) GetKey() string {
	return setting.key
}

func (setting Setting) GetValue() string {
	return setting.value
}

func (setting Setting) ToString() string {
	return setting.key + " = " + setting.value
}

type SettingList []Setting

func AppendEnvList(settingList SettingList, envs ...Setting) SettingList {
	return append(settingList, envs...)
}

func (settingList SettingList) ToString() string {
	temp := make([]string, len(settingList))
	for i, env := range settingList {
		temp[i] = env.ToString()
	}
	return strings.Join(temp, "\n")
}
