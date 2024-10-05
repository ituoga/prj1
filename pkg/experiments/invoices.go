package experiments

import (
	"context"
	"log"

	"github.com/delaneyj/toolbelt"
	"github.com/google/uuid"
	types "github.com/ituoga/prj1/types"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type InvoiceDB struct {
	db *toolbelt.Database
}

func Invoice() *InvoiceDB {
	return &InvoiceDB{db: db}
}

func (i *InvoiceDB) Store(ctx context.Context, inv types.Invoice) error {
	uid := uuid.NewString()
	err := i.db.WriteTX(ctx, func(tx *sqlite.Conn) error {
		return sqlitex.Execute(tx,
			"insert into invoices (uuid, document_date, due_date, currency, rate, serial_name, recipient_name, recipient_code, recipient_vat, recipient_email, recipient_phone, recipient_addr, recipient_country, written_by) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?)", &sqlitex.ExecOptions{
				Args: []interface{}{
					uid,
					inv.DocumentDate,
					inv.DueDate,
					inv.Currency,
					inv.Rate,
					inv.SerialName,
					inv.RecipientName,
					inv.RecipientCode,
					inv.RecipientVAT,
					inv.RecipientEmail,
					inv.RecipientPhone,
					inv.RecipientAddr,
					inv.RecipientCountry,
					inv.WrittenBy,
				},
			})
	})
	if err != nil {
		return err
	}
	err = i.db.WriteTX(ctx, func(tx *sqlite.Conn) error {
		for _, row := range inv.Lines {
			return sqlitex.Execute(tx,
				"insert into invoice_rows (invoice_id, name, price, comment, qty, units, vat, uid) values (?,?,?,?,?,?,?,?)", &sqlitex.ExecOptions{
					Args: []interface{}{
						uid,
						row.Name,
						row.Price,
						row.Comment,
						row.Quantity,
						row.Units,
						row.Vat,
						row.UID,
					},
				})
		}
		return nil
	})

	return err
}

func (i *InvoiceDB) LoadRecipient(ctx context.Context, id string, iv *types.Invoice) error {
	var err error
	err = i.db.ReadTX(ctx, func(tx *sqlite.Conn) error {
		return sqlitex.Execute(tx, "select * from invoices where id = ? limit 1", &sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				iv.RecipientName = stmt.ColumnText(7)
				iv.RecipientCode = stmt.ColumnText(8)
				iv.RecipientVAT = stmt.ColumnText(9)
				iv.RecipientEmail = stmt.ColumnText(10)
				iv.RecipientPhone = stmt.ColumnText(11)
				iv.RecipientAddr = stmt.ColumnText(12)
				iv.RecipientCountry = stmt.ColumnText(13)
				return nil
			},
			Args: []interface{}{id},
		})
	})
	return err
}

func (i *InvoiceDB) Load(ctx context.Context, id string) (types.Store[types.Invoice], error) {
	var err error
	var inv types.Store[types.Invoice]
	err = i.db.ReadTX(ctx, func(tx *sqlite.Conn) error {
		return sqlitex.Execute(tx, "select * from invoices where id = ? limit 1", &sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				*inv.UUID = stmt.ColumnText(1)
				inv.Data.DocumentDate = stmt.ColumnText(2)
				inv.Data.DueDate = stmt.ColumnText(3)
				inv.Data.Currency = stmt.ColumnText(4)
				inv.Data.Rate = stmt.ColumnFloat(5)
				inv.Data.SerialName = stmt.ColumnText(6)
				inv.Data.RecipientName = stmt.ColumnText(7)
				inv.Data.RecipientCode = stmt.ColumnText(8)
				inv.Data.RecipientVAT = stmt.ColumnText(9)
				inv.Data.RecipientEmail = stmt.ColumnText(10)
				inv.Data.RecipientPhone = stmt.ColumnText(11)
				inv.Data.RecipientAddr = stmt.ColumnText(12)
				inv.Data.RecipientCountry = stmt.ColumnText(13)
				inv.Data.WrittenBy = stmt.ColumnText(14)
				return nil
			},
			Args: []interface{}{id},
		})
	})
	if err != nil {
		return inv, err
	}

	var rows []types.InvoiceRow
	err = i.db.ReadTX(ctx, func(tx *sqlite.Conn) error {
		return sqlitex.Execute(tx, "select * from invoice_rows where invoice_id = ?", &sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				var row types.InvoiceRow
				row.UID = stmt.ColumnText(1)
				row.Name = stmt.ColumnText(2)
				row.Price = stmt.ColumnFloat(3)
				row.Comment = stmt.ColumnText(4)
				row.Quantity = stmt.ColumnFloat(5)
				row.Units = stmt.ColumnText(6)
				row.Vat = stmt.ColumnFloat(7)
				rows = append(rows, row)
				return nil
			},
			Args: []interface{}{id},
		})
	})

	inv.Data.Lines = append(inv.Data.Lines, rows...)

	return inv, err
}

func (i *InvoiceDB) List(ctx context.Context) ([]types.Store[types.Invoice], error) {
	var err error
	var invs []types.Store[types.Invoice]
	err = i.db.ReadTX(ctx, func(tx *sqlite.Conn) error {
		return sqlitex.Execute(tx, "select * from invoices", &sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				var inv types.Store[types.Invoice]
				uid := stmt.ColumnText(1)
				inv.UUID = &uid
				inv.Data.DocumentDate = stmt.ColumnText(2)
				inv.Data.DueDate = stmt.ColumnText(3)
				inv.Data.Currency = stmt.ColumnText(4)
				inv.Data.Rate = stmt.ColumnFloat(5)
				inv.Data.SerialName = stmt.ColumnText(6)
				inv.Data.RecipientName = stmt.ColumnText(7)
				inv.Data.RecipientCode = stmt.ColumnText(8)
				inv.Data.RecipientVAT = stmt.ColumnText(9)
				inv.Data.RecipientEmail = stmt.ColumnText(10)
				inv.Data.RecipientPhone = stmt.ColumnText(11)
				inv.Data.RecipientAddr = stmt.ColumnText(12)
				inv.Data.RecipientCountry = stmt.ColumnText(13)
				inv.Data.WrittenBy = stmt.ColumnText(14)
				invs = append(invs, inv)
				return nil
			},
		})
	})
	if err != nil {
		log.Printf("%v", err)
	}
	return invs, nil
}
