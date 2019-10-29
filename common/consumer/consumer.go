// Package consumer provides basic interfaces for consumer
package consumer

import (
	"io"
	"strconv"

	"github.com/alexbyk/goquiz/common/model"
)

// Writer is the interface that allows to save customers in chunks
type Writer interface {
	WriteCustomers(customers []*model.Customer) (int, error)
}

// Notifier is the interface to notify that something happened
type Notifier interface {
	Notify(string) error
}

// Reader is the interface that wraps ReadCustomers method
type Reader interface {

	// ReadCustomers should read records in chunks. It returns io.EOF (maybe with last records) in the end
	// If any error occurs, it also returns records that could parse
	ReadCustomers() ([]*model.Customer, error)
}

// Consumer implement consumer logic
type Consumer struct {
	writer   Writer
	reader   Reader
	notifier Notifier
}

// NewConsumer returns a new consumer
func NewConsumer(r Reader, w Writer, n Notifier) *Consumer {
	return &Consumer{reader: r, writer: w, notifier: n}
}

// Consume consumes an input
func (c *Consumer) Consume() error {
	var eof bool
	for !eof {
		customers, err := c.reader.ReadCustomers()
		if err == io.EOF {
			eof = true
		} else if err != nil {
			return err
		}
		n, err := c.writer.WriteCustomers(customers)
		if err != nil {
			return err
		}
		err = c.notifier.Notify(strconv.Itoa(n))
	}
	return nil
}
