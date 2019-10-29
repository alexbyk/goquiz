// Package consumer provides basic interfaces for consumer
package consumer

import "github.com/alexbyk/goquiz/model"

// Writer is the interface that allows to save customers in chunks
type Writer interface {
	WriteRecords(customers []model.Customer) error
}

// Reader is the interface that wraps ReadRecords method
type Reader interface {

	// ReadRecords should read records in chunks. It returns io.EOF (maybe with last records) in the end
	// If any error occurs, it also returns records that could parse
	ReadRecords() ([]*model.Customer, error)
}
