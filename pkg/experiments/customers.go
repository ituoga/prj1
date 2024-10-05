package experiments

import (
	"context"

	"github.com/delaneyj/toolbelt"
	"github.com/google/uuid"
	types "github.com/ituoga/prj1/types"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type CustomersDB struct {
	db *toolbelt.Database
}

func Customer() *CustomersDB {
	return &CustomersDB{db: db}
}

func (cdb *CustomersDB) Store(c types.Customer, luuid *string) error {
	if luuid == nil {
		var s = uuid.NewString()
		luuid = &s
	}
	if *luuid == "" {
		var s = uuid.NewString()
		luuid = &s
	}
	err := cdb.db.WriteTX(context.TODO(), func(tx *sqlite.Conn) error {
		return sqlitex.ExecuteTransient(tx, "insert into customers (uuid, name, code, vat, email, phone, address, country) values (?, ?, ?, ?, ?, ?, ?, ?)", &sqlitex.ExecOptions{
			Args: []interface{}{*luuid, c.Name, c.Code, c.VAT, c.Email, c.Phone, c.Addr, c.Country},
		})
	})
	if err != nil {
		return err
	}
	return nil
}

func (cdb *CustomersDB) Update(c types.CustomerStore) error {
	err := cdb.db.WriteTX(context.TODO(), func(tx *sqlite.Conn) error {
		return sqlitex.ExecuteTransient(tx, "update customers set name = ?, code = ?, vat = ?, email = ?, phone = ?, address = ?, country = ? where uuid = ?", &sqlitex.ExecOptions{
			Args: []interface{}{c.Data.Name, c.Data.Code, c.Data.VAT, c.Data.Email, c.Data.Phone, c.Data.Addr, c.Data.Country, c.UUID},
		})
	})
	if err != nil {
		return err
	}
	return nil
}

func (cdb *CustomersDB) Load(uuid string) types.CustomerStore {
	var c types.CustomerStore
	err := cdb.db.ReadTX(context.TODO(), func(tx *sqlite.Conn) error {
		return sqlitex.Execute(tx, "select id, uuid, name, code, vat, email, phone, address, country from customers where uuid = ? limit 1", &sqlitex.ExecOptions{
			Args: []interface{}{uuid},
			ResultFunc: func(stmt *sqlite.Stmt) error {
				c.ID = stmt.ColumnInt(0)
				uid := stmt.ColumnText(1)
				c.UUID = &uid
				c.Data.Name = stmt.ColumnText(2)
				c.Data.Code = stmt.ColumnText(3)
				c.Data.VAT = stmt.ColumnText(4)
				c.Data.Email = stmt.ColumnText(5)
				c.Data.Phone = stmt.ColumnText(6)
				c.Data.Addr = stmt.ColumnText(7)
				c.Data.Country = stmt.ColumnText(8)
				return nil
			},
		})
	})
	if err != nil {
		return types.CustomerStore{}
	}
	return c
}

func (cdb *CustomersDB) List() []types.CustomerStore {
	var cs []types.CustomerStore
	err := cdb.db.ReadTX(context.TODO(), func(tx *sqlite.Conn) error {
		return sqlitex.Execute(tx, "select id, uuid, name, code, vat, email, phone, address, country from customers", &sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				var tcs types.CustomerStore
				tcs.ID = stmt.ColumnInt(0)
				uid := stmt.ColumnText(1)
				tcs.UUID = &uid
				tcs.Data.Name = stmt.ColumnText(2)
				tcs.Data.Code = stmt.ColumnText(3)
				tcs.Data.VAT = stmt.ColumnText(4)
				tcs.Data.Email = stmt.ColumnText(5)
				tcs.Data.Phone = stmt.ColumnText(6)
				tcs.Data.Addr = stmt.ColumnText(7)
				tcs.Data.Country = stmt.ColumnText(8)
				cs = append(cs, tcs)
				return nil
			},
		})
	})
	if err != nil {
		return nil
	}
	return cs
}

func (cdb *CustomersDB) Delete(uuid string) error {
	err := cdb.db.WriteTX(context.TODO(), func(tx *sqlite.Conn) error {
		return sqlitex.ExecuteTransient(tx, "delete from customers where uuid = ?", &sqlitex.ExecOptions{
			Args: []interface{}{uuid},
		})
	})
	if err != nil {
		return err
	}
	return nil
}
