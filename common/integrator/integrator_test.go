package integrator_test

import (
	"testing"
	"time"

	"github.com/alexbyk/ftest"
	"github.com/alexbyk/goquiz/common/integrator"

	"github.com/alexbyk/goquiz/common/model"
)

// Mocks
type mNexter struct {
	data []*model.Customer
}

func (m *mNexter) Next(fn func(*model.Customer) error) error {
	if len(m.data) == 0 {
		return integrator.ErrEmpty
	}
	var c *model.Customer
	c = m.data[0]
	err := fn(c)
	if err == nil {
		m.data = m.data[1:]
	}
	return err
}

type mListener struct {
	c chan string
}

func (l *mListener) Listen() <-chan string {
	l.c = make(chan string)
	return l.c
}

func (l *mListener) Unlisten() error {
	close(l.c)
	return nil
}

type mPublisher struct {
	data []error
	got  []*model.Customer
}

func (p *mPublisher) Publish(c *model.Customer) error {
	p.got = append(p.got, c)
	if len(p.data) == 0 {
		return nil
	}
	var e error
	e, p.data = p.data[0], p.data[1:]
	return e
}

type mLogger struct{}

func (*mLogger) Debug(...interface{}) {}

func TestUsage2(t *testing.T) {
	ft := ftest.New(t)
	l := new(mListener)
	n := new(mNexter)
	p := new(mPublisher)
	intg := integrator.NewIntegrator(n, l, p)
	intg.Logger = &mLogger{}
	data := []*model.Customer{{ID: "1"}, {ID: "2"}}

	n.data = data
	go func() {
		time.Sleep(time.Millisecond * 300)
		intg.Stop()
	}()
	intg.Start()

	ft.Eq(p.got, data)
}
