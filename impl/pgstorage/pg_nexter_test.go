package pgstorage_test

import (
	"errors"
	"testing"

	"github.com/alexbyk/goquiz/common/integrator"

	"github.com/alexbyk/goquiz/common/model"

	"github.com/alexbyk/ftest"
	"github.com/alexbyk/goquiz/impl/pgstorage"
)

func TestNext(t *testing.T) {
	db := skipDSN(t)
	ft := ftest.New(t)
	n := pgstorage.NewPgNexter(db)
	err := db.Insert(&pgstorage.Customer{Customer: &model.Customer{ID: "1"}})
	err = db.Insert(&pgstorage.Customer{Customer: &model.Customer{ID: "2"}})
	ft.Nil(err)

	checkStatus := func(id string, expect string) {
		customer := &pgstorage.Customer{}
		db.Model(customer).Where(`id = ?`, id).Select()
		ft.Eq(customer.Status, expect)
	}

	err = n.Next(func(c *model.Customer) error {
		ft.Eq(c.ID, "1")
		checkStatus("1", pgstorage.CustomerStatusSending)
		checkStatus("2", pgstorage.CustomerStatusEmpty)
		return nil
	})
	ft.Nil(err)
	checkStatus("1", pgstorage.CustomerStatusConfirmed)
	checkStatus("2", pgstorage.CustomerStatusEmpty)

	// fail second
	n.Next(func(c *model.Customer) error { return errors.New("Mock") })
	checkStatus("1", pgstorage.CustomerStatusConfirmed)
	checkStatus("2", pgstorage.CustomerStatusEmpty)

	// ok second
	n.Next(func(c *model.Customer) error { return nil })
	checkStatus("1", pgstorage.CustomerStatusConfirmed)
	checkStatus("2", pgstorage.CustomerStatusConfirmed)

	// no records
	err = n.Next(func(c *model.Customer) error { return nil })
	ft.Eq(err, integrator.ErrEmpty)
}
