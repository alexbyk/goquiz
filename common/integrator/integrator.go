package integrator

import (
	"errors"

	"github.com/alexbyk/goquiz/common/model"
)

// ErrEmpty is a flag that no records left
var ErrEmpty = errors.New("No customers")

// Nexter reads customers one by one
type Nexter interface {

	// Next handle next records. It's up to implementation detail how to deal with errors. Should return ErrEmpty if no records left
	// Passed function will be executed until it returns nil
	Next(func(*model.Customer) error) error
}

// Listener reacts if there is something to read
type Listener interface {
	Listen() <-chan string
	Unlisten() error
}

// Publisher publish customers one by one to 3d party service
type Publisher interface {
	Publish(*model.Customer) error
}
