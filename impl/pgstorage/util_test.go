package pgstorage_test

import (
	"os"
	"testing"

	"github.com/alexbyk/goquiz/impl/pgstorage"
	"github.com/go-pg/pg"
)

func skipDSN(t *testing.T) *pg.DB {
	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		t.Skip("Provide TEST_DSN env")
	}
	db, err := pgstorage.Connect(dsn)
	if err != nil {
		t.Fatal(err)
	}
	if err != pgstorage.CreateTable(db) {
		t.Fatal(err)
	}
	if err != pgstorage.Truncate(db) {
		t.Fatal(err)
	}
	return db
}
