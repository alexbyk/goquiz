package pgstorage_test

import (
	"fmt"
	"testing"

	"github.com/alexbyk/goquiz/impl/pgstorage"

	"github.com/alexbyk/ftest"
)

func TestListener(t *testing.T) {
	ft := ftest.New(t)
	db := skipDSN(t)
	l := pgstorage.NewPgListener(db)
	ch := l.Listen()
	go func() {
		_, err := db.Exec(fmt.Sprintf("NOTIFY \"%s\", ?", pgstorage.Channel), "Foo")
		ft.Nil(err)
	}()
	ft.Eq(<-ch, "Foo")
}
