package experiments

import (
	"context"
	"errors"
	"log"

	"github.com/delaneyj/toolbelt"
	"github.com/google/uuid"
	"github.com/ituoga/proj1/pkg/password"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

var db *toolbelt.Database

func Auth(user, pass string) bool {
	if len(user) == 0 || len(pass) == 0 {
		return false
	}
	found := false
	ctx := context.TODO()
	err := db.ReadTX(ctx, func(tx *sqlite.Conn) error {
		return sqlitex.Execute(tx, "select password from users where username = ? limit 1", &sqlitex.ExecOptions{
			Args: []interface{}{user},
			ResultFunc: func(stmt *sqlite.Stmt) error {
				log.Printf("=> %v %v %v %v", stmt.ColumnText(0), user, pass, password.CheckHash(pass, stmt.ColumnText(0)))
				if !password.CheckHash(pass, stmt.ColumnText(0)) {
					return errors.New("password mismatch")
				}
				found = true
				return nil
			},
		})
	})
	if err != nil {
		log.Printf("%v", err)
	}
	return found
}

func DB() {
	ctx := context.TODO()
	var err error
	db, err = toolbelt.NewDatabase(ctx, "./data/db.sqlite", []string{
		"CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, uuid varchar(255), username varchar(255), password varchar(255), created_at datetime, updated_at datetime, confirmed_at datetime)",
	})
	if err != nil {
		log.Printf("database error: %v", err)
		return
	}

	h, _ := password.Hash("admin")

	createRecord := func(in string) error {
		err = db.WriteTX(ctx, func(tx *sqlite.Conn) error {
			err = sqlitex.ExecuteTransient(tx, "insert into users (uuid, username, password) values (?,?,?)", &sqlitex.ExecOptions{
				Args: []interface{}{uuid.NewString(), "admin", h},
			})
			if err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			log.Printf("database error: %v", err)
			return err
		}
		return nil
	}
	_ = createRecord
	// log.Fatal(createRecord("admin"))

	// d.ReadTX(ctx, func(tx *sqlite.Conn) error {
	// 	sqlitex.Execute(tx, "select * from projects", &sqlitex.ExecOptions{
	// 		ResultFunc: func(stmt *sqlite.Stmt) error {
	// 			log.Printf("project: %v", stmt.ColumnText(1))
	// 			return nil
	// 		},
	// 	})
	// 	return nil
	// })

}
