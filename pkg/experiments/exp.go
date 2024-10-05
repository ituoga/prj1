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

	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var db *toolbelt.Database
var xdb *sqlx.DB

func init() {
	DB()

	dba, err := sql.Open("sqlite", "__deleteme.db")
	if err != nil {
		panic(err)
	}
	xdb = sqlx.NewDb(dba, "sqlite")

	for _, m := range Migrations() {
		tx, err := xdb.Begin()
		if err != nil {
			log.Printf("transaction error: %v", err)
		}
		_, err = tx.Exec(m)
		if err != nil {
			log.Printf("migration error: %v", err)
			tx.Rollback()
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Printf("commit error: %v", err)
		}
	}
}

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

func Migrations() []string {
	return []string{
		"CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, uuid varchar(255), username varchar(255), password varchar(255), created_at datetime, updated_at datetime, confirmed_at datetime)",
		"CREATE TABLE IF NOT EXISTS settings (id INTEGER PRIMARY KEY, uuid varchar(255), name varchar(255), code varchar(255), vat varchar(255), email varchar(255), phone varchar(255), address varchar(255), country varchar(255), serial_name varchar(255), serial_no int, created_at datetime, updated_at datetime)",
		"CREATE TABLE IF NOT EXISTS invoices (id INTEGER PRIMARY KEY, uuid varchar(255), document_date datetime, due_date datetime, currency varchar(255), rate float, serial_name varchar(255), recipient_name varchar(255), recipient_code varchar(255), recipient_vat varchar(255), recipient_email varchar(255), recipient_phone varchar(255), recipient_addr varchar(255), recipient_country varchar(255), written_by varchar(255), taken_by varchar(255), created_at datetime, updated_at datetime)",
		"CREATE TABLE IF NOT EXISTS invoice_rows (id INTEGER PRIMARY KEY, invoice_id INTEGER, number int, name varchar(255), price float, comment varchar(255), uid varchar(255), qty float, units varchar(255), vat float, created_at datetime, updated_at datetime)",
		"CREATE TABLE IF NOT EXISTS customers (id INTEGER PRIMARY KEY, uuid varchar(255), name varchar(255), code varchar(255), vat varchar(255), email varchar(255), phone varchar(255), address varchar(255), country varchar(255), created_at datetime, updated_at datetime)",
		"CREATE TABLE IF NOT EXISTS products (id INTEGER PRIMARY KEY, uuid varchar(255), name varchar(255), code varchar(255), created_at datetime, updated_at datetime)",
	}
}

func DB() {
	if db != nil || xdb != nil {
		return
	}
	ctx := context.TODO()
	var err error
	db, err = toolbelt.NewDatabase(ctx, "./data/db.sqlite", Migrations())
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
}
