package experiments

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/ituoga/proj1/types"
)

type AutoCompleteDB struct {
	db *sqlx.DB
}

func AutoComplete() *AutoCompleteDB {
	return &AutoCompleteDB{db: xdb}
}

func (a *AutoCompleteDB) Get(ctx context.Context, id string) (types.CustomerStore, error) {
	var results types.CustomerStore

	row := a.db.QueryRowxContext(ctx, "select id,name,code,vat,phone,email,address,country from customers where id = ?", id)
	err := row.Scan(&results.ID, &results.Data.Name, &results.Data.Code, &results.Data.VAT, &results.Data.Phone, &results.Data.Email, &results.Data.Addr, &results.Data.Country)

	return results, err
}

func (a *AutoCompleteDB) List(q string) ([]types.Complete, error) {
	var results []types.Complete
	err := a.db.Select(&results, "select id,name from customers where name like ?", "%"+q+"%")
	if err != nil {
		return nil, err
	}
	return results, nil
}
