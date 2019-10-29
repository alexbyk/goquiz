package pgstorage

import (
	"github.com/alexbyk/goquiz/common/integrator"
	"github.com/alexbyk/goquiz/common/model"
	"github.com/go-pg/pg"
)

// PgNexter implements integrator.Nexter interface
type PgNexter struct {
	db *pg.DB
}

// NewPgNexter returns PgNexter
func NewPgNexter(db *pg.DB) *PgNexter {
	return &PgNexter{db: db}
}

// fetch next record and mark as "sending"
func (p *PgNexter) fetchNext() (*Customer, error) {
	rec := &Customer{}
	err := p.db.RunInTransaction(func(tx *pg.Tx) error {
		if err := tx.Model(rec).Where("status = ?", CustomerStatusEmpty).For("Update").First(); err != nil {
			return err
		}
		rec.Status = CustomerStatusSending
		_, err := tx.Model(rec).WherePK().Column("status").Returning("*").Update()
		return err
	})
	if err != nil {
		return nil, err
	}
	return rec, nil
}

// Next process one record. No transaction required because integrator suppose to handle this
func (p *PgNexter) Next(fn func(*model.Customer) error) error {

	customer, err := p.fetchNext()

	if err == pg.ErrNoRows {
		return integrator.ErrEmpty
	} else if err != nil {
		return err
	}
	if err := fn(customer.Customer); err != nil {
		customer.Status = CustomerStatusEmpty
		if _, err := p.db.Model(customer).WherePK().Column("status").Update(); err != nil {
			return err
		}
		return err
	}

	customer.Status = CustomerStatusConfirmed
	if _, err := p.db.Model(customer).WherePK().Column("status").Update(); err != nil {
		return err
	}
	return nil

}
