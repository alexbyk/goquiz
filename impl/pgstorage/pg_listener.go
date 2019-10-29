package pgstorage

import (
	"sync"

	"github.com/go-pg/pg"
)

// PgListener implements Listener backend
type PgListener struct {
	db       *pg.DB
	mu       sync.Mutex
	listener *pg.Listener
}

// NewPgListener returns a new PgListener but dosn't listen anything
func NewPgListener(db *pg.DB) *PgListener {
	return &PgListener{db: db}
}

func (p *PgListener) ensureListener() *pg.Listener {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.listener == nil {
		p.listener = p.db.Listen(Channel)
	}
	return p.listener
}

// Listen returns a channel
func (p *PgListener) Listen() <-chan string {
	p.ensureListener()
	chStr := make(chan string)
	ch := p.ensureListener().Channel()
	go func() {
		for v := range ch {
			chStr <- v.Payload
		}
		close(chStr)
	}()
	return chStr
}

// Unlisten clears resources and closes channel
func (p *PgListener) Unlisten() error {
	p.mu.Lock()
	l := p.listener
	if l == nil {
		return nil
	}
	p.listener = nil
	p.mu.Unlock()
	return l.Close()
}
