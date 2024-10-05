package pages

import types "github.com/ituoga/prj1/types"

type ProductsIndex struct {
	DS   types.DataStore
	Data []types.ProductStore
}

type ProductsForm struct {
	DS   types.DataStore
	Form types.ProductStore
}
