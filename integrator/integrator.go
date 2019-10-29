package integrator

import "github.com/alexbyk/goquiz/model"

// Nexter reads customers one by one
type Nexter interface {
	Next(*model.Customer) error
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

// ListenNexter combines Listener and Nexter interfaces
type ListenNexter interface {
	Listener
	Nexter
}
