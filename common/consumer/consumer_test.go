package consumer_test

import (
	"io"
	"strconv"
	"testing"

	"github.com/alexbyk/ftest"
	"github.com/alexbyk/goquiz/common/consumer"
	"github.com/alexbyk/goquiz/common/model"
)

// mock implementation
type mReader struct {
	id int
}

func (m *mReader) ReadCustomers() ([]*model.Customer, error) {
	m.id++
	var err error
	if m.id > 2 {
		err = io.EOF
	}
	return []*model.Customer{{ID: strconv.Itoa(m.id)}}, err
}

type mWriter struct {
	got []string
}

func (m *mWriter) WriteCustomers(data []*model.Customer) (int, error) {
	for _, r := range data {
		m.got = append(m.got, r.ID)
	}
	return len(data), nil
}

type mNotifier struct {
	got []string
}

func (m *mNotifier) Notify(payload string) error {
	m.got = append(m.got, payload)
	return nil
}

func TestConsume(t *testing.T) {
	ft := ftest.New(t)
	r := new(mReader)
	w := new(mWriter)
	n := new(mNotifier)
	c := consumer.NewConsumer(r, w, n)
	err := c.Consume()
	ft.Nil(err)
	ft.Eq(w.got, []string{"1", "2", "3"})
	ft.Eq(n.got, []string{"1", "1", "1"})
}
