package experiments

import (
	"context"
	"log"
	"strconv"

	"github.com/delaneyj/toolbelt"
	types "github.com/ituoga/prj1/types"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type SettingsDB struct {
	db *toolbelt.Database
}

func Settings() *SettingsDB {
	return &SettingsDB{db: db}
}

func (s *SettingsDB) Store(ctx context.Context, settings types.Settings) {
	err := s.db.WriteTX(ctx, func(tx *sqlite.Conn) error {
		err := sqlitex.Execute(tx, "delete from settings where id = ?", &sqlitex.ExecOptions{
			Args: []interface{}{"1"},
		})
		if err != nil {
			log.Printf("%v", err)
			return err
		}
		return sqlitex.ExecuteTransient(tx, `
			insert into settings (id,name, code, vat, email, phone, address, country, serial_name, serial_no, created_at, updated_at)
			values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
		`, &sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				return nil
			},
			Args: []interface{}{
				"1", // user id
				settings.MyName,
				settings.MyCode,
				settings.MyVAT,
				settings.MyEmail,
				settings.MyPhone,
				settings.MyAddr,
				settings.MyCountry,
				settings.SerialName,
				settings.SerialNo,
			},
		})
	})
	if err != nil {
		log.Printf("===== %v", err)
	}
}

func (s *SettingsDB) Load(ctx context.Context) types.Settings {
	var results = &types.Settings{}
	err := s.db.ReadTX(ctx, func(tx *sqlite.Conn) error {
		return sqlitex.Execute(tx, "select * from settings where id = ?", &sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				results.MyName = stmt.ColumnText(2)
				results.MyCode = stmt.ColumnText(3)
				results.MyVAT = stmt.ColumnText(4)
				results.MyEmail = stmt.ColumnText(5)
				results.MyPhone = stmt.ColumnText(6)
				results.MyAddr = stmt.ColumnText(7)
				results.MyCountry = stmt.ColumnText(8)
				results.SerialName = stmt.ColumnText(9)
				results.SerialNo, _ = strconv.Atoi(stmt.ColumnText(10))
				return nil
			},
			Args: []interface{}{"1"},
		})
	})
	if err != nil {
		log.Printf("%v", err)
	}
	return *results
}
