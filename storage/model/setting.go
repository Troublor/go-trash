package model

import "time"

// TODO generate default setting configuration json file
//go:generate

type Setting struct {
	AutoRemove         bool          `json:"autoRemove"`         // whether to auto remove trashes in the bin
	AutoRemoveInterval time.Duration `json:"autoRemoveInterval"` // only valid when AutoRemove == true, auto remove interval period
}
