package pgstorage

import (
	"fmt"

	"github.com/go-pg/pg"
)

// PgNotifier implement consumer.Notifier interface in Postgres
type PgNotifier struct {
	db *pg.DB
}

// NewPgNotifier returns new notifier
func NewPgNotifier(db *pg.DB) *PgNotifier {
	return &PgNotifier{db: db}
}

// Notify sends a notification
func (p *PgNotifier) Notify(payload string) error {
	_, err := p.db.Exec(fmt.Sprintf(`NOTIFY "%s", ?`, Channel), payload)
	return err
}
