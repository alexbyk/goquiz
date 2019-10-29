package pgstorage_test

import (
	"testing"

	"github.com/alexbyk/goquiz/common/consumer"

	"github.com/alexbyk/ftest"
	"github.com/alexbyk/goquiz/common/model"
	"github.com/alexbyk/goquiz/impl/pgstorage"
)

func TestWr(t *testing.T) {
	db := skipDSN(t)
	var wr consumer.Writer = pgstorage.NewPgWriter(db)
	ft := ftest.New(t)
	count, err := wr.WriteCustomers([]*model.Customer{{ID: "1"}, {ID: "1"}})
	ft.Nil(err).Eq(count, 1)
	customer := &model.Customer{}

	db.Model(customer).Where("id = ?", "1").Select()
	ft.Eq(customer.ID, "1")
}
