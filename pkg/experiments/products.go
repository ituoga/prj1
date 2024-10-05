package experiments

import (
	"log"

	types "github.com/ituoga/prj1/types"
	"github.com/jmoiron/sqlx"
)

type ProductDB struct {
	db *sqlx.DB
}

func Product() *ProductDB {
	return &ProductDB{xdb}
}

func (p *ProductDB) List() ([]types.ProductStore, error) {
	var projects []types.ProductStore

	rows, err := p.db.Queryx("SELECT id,uuid,name,code FROM products")

	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	for rows.Next() {
		var project types.ProductStore
		err = rows.Scan(&project.ID, &project.UUID, &project.Data.Name, &project.Data.Code)
		if err != nil {
			log.Printf("error: %v", err)
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil

}

func (p *ProductDB) Get(id string) (types.ProductStore, error) {
	var project types.ProductStore

	row := p.db.QueryRowx("SELECT id,uuid,name,code FROM products WHERE uuid = ?", id)
	err := row.Scan(&project.ID, &project.UUID, &project.Data.Name, &project.Data.Code)

	return project, err
}

func (p *ProductDB) Create(project *types.ProductStore) error {
	_, err := p.db.NamedExec("INSERT INTO products (uuid, name, code) VALUES (:uuid, :name, :code)", project.ToMap())
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductDB) Update(project *types.ProductStore) error {
	_, err := p.db.NamedExec("UPDATE products SET name = :name, code = :code WHERE uuid = :uuid", project.ToMap())
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductDB) Delete(id string) error {
	_, err := p.db.Exec("DELETE FROM products WHERE uuid = ?", id)
	if err != nil {
		return err
	}
	return nil
}
