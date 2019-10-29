package pgstorage_test

import (
	"testing"

	"github.com/alexbyk/goquiz/impl/pgstorage"

	"github.com/alexbyk/ftest"
)

func TestNotifier(t *testing.T) {
	ft := ftest.New(t)
	db := skipDSN(t)
	l := db.Listen(pgstorage.Channel)
	defer l.Close()

	n := pgstorage.NewPgNotifier(db)
	go func() {
		n.Notify("hello")
	}()
	ft.Eq((<-l.Channel()).Payload, "hello")
}
