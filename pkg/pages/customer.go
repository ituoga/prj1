package pages

import (
	types "github.com/ituoga/prj1/types"
)

type CustomerForm struct {
	DS   types.DataStore
	Form types.CustomerStore

	Settings *types.Settings
}

type CustomerIndex struct {
	DS   types.DataStore
	Data []types.CustomerStore
}
