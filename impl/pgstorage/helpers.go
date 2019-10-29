package pgstorage

import (
	"github.com/alexbyk/goquiz/common/model"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Channel is a name of the channel
const Channel = "app:customers"

// Truncate clears a "customers" table for testing purposes
func Truncate(db *pg.DB) error {
	_, err := db.Exec("TRUNCATE TABLE customers RESTART IDENTITY CASCADE")
	return err
}

// Connect returns a pg.DB instance, connected by parsed dsn
func Connect(dsn string) (*pg.DB, error) {
	opts, err := pg.ParseURL(dsn)
	if err != nil {
		return nil, err
	}

	db := pg.Connect(opts)
	return db, nil
}

// CreateTable creates a customers table if neccesarry
func CreateTable(db *pg.DB) error {
	return db.CreateTable((*model.Customer)(nil), &orm.CreateTableOptions{IfNotExists: true})
}
