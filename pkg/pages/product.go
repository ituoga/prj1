package pages

import "github.com/ituoga/proj1/types"

type ProductsIndex struct {
	DS   types.DataStore
	Data []types.ProductStore
}

type ProductsForm struct {
	DS   types.DataStore
	Form types.ProductStore
}
