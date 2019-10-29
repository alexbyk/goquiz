/*Package pgstorage provides Postgres based implementation for common interfaces*/
package pgstorage

import (
	"github.com/alexbyk/goquiz/common/model"
	"github.com/go-pg/pg"
)

// PgWriter implements Postgres consumer.Writer backend
type PgWriter struct {
	db *pg.DB
}

// NewPgWriter returns PgWriter
func NewPgWriter(db *pg.DB) *PgWriter {
	return &PgWriter{db}
}

// WriteCustomers saves records to db, ignoring duplicates
func (s *PgWriter) WriteCustomers(records []*model.Customer) (int, error) {
	m := ConvertModels(records)
	ret, err := s.db.Model(&m).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return 0, err
	}
	return ret.RowsAffected(), nil
}
