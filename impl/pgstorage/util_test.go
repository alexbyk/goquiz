package pgstorage_test

import (
	"os"
	"testing"

	"github.com/alexbyk/goquiz/impl/pgstorage"
	"github.com/go-pg/pg"
)

// docker run --rm -p 5432:5432 -e POSTGRES_USER=test  postgres:alpine postgres -c log_statement=all
// TEST_DSN="postgres://test@localhost/test?sslmode=disable" go test -count=1 -v ./...  -timeout 2s

func skipDSN(t *testing.T) *pg.DB {
	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		t.Skip("Provide TEST_DSN env")
	}
	db, err := pgstorage.Connect(dsn)
	if err != nil {
		t.Fatal(err)
	}
	if err = pgstorage.CreateTable(db); err != nil {
		t.Fatal(err)
	}
	if err = pgstorage.Truncate(db); err != nil {
		t.Fatal(err)
	}
	return db
}
