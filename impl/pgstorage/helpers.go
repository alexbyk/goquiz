package pgstorage

import (
	"github.com/alexbyk/goquiz/common/model"
	"github.com/go-pg/pg"
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
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS "customers" (
  id bigint NOT NULL,
  first_name text NOT NULL default '',
  last_name text NOT NULL default '',
  email text NOT NULL default '',
  phone text NOT NULL default '',
  status text NOT NULL default '',
  PRIMARY KEY ("id"))`)
	return err
}

// ConvertModel converts general Customer to our DB Customer
func ConvertModel(c *model.Customer) *Customer {
	return &Customer{Customer: c}
}

// ConvertModels converts array of models
func ConvertModels(arr []*model.Customer) []*Customer {
	ret := []*Customer{}
	for _, cur := range arr {
		ret = append(ret, ConvertModel(cur))
	}
	return ret
}
