package integrator

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/alexbyk/goquiz/common/model"
)

// DefaultDelayRetry is a default value for Integrator.DelayRetry
const DefaultDelayRetry = time.Second * 10

// DefaultDelayTick is a default value for Integrator.DelayTick
const DefaultDelayTick = time.Second * 4

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

// Logger is a simple logger interface
type Logger interface {
	Debug(...interface{})
}

// SimpleLogger is a dummy logger
type SimpleLogger struct{}

// Debug prints everything to the output
func (logger *SimpleLogger) Debug(args ...interface{}) { log.Println(args...) }

// Integrator sends data to the 3d party
type Integrator struct {
	nexter    Nexter
	listener  Listener
	publisher Publisher
	Logger    Logger

	// Delay between attempt to retry Publish
	DelayRetry time.Duration

	chExit chan struct{} // channel and a flag that we started

	muRunning sync.Mutex
	running   bool
}

// NewIntegrator returns an Integrator instance
func NewIntegrator(n Nexter, l Listener, p Publisher) *Integrator {
	return &Integrator{nexter: n, listener: l, publisher: p, DelayRetry: DefaultDelayRetry, Logger: &SimpleLogger{}}
}

// Start starts a daemon loop, sholdn't be called from multiple threads
func (i *Integrator) Start() {

	chListen := i.listener.Listen()
	defer i.listener.Unlisten()

	i.chExit = make(chan struct{})
	defer close(i.chExit)

	i.maybeRun() // run at the beginning
	for {
		select {
		case <-i.chExit:
			log.Println("exit")
			return
		case payload := <-chListen:
			log.Println("Payload: ", payload)
		}
		i.maybeRun()
	}
}

func (i *Integrator) maybeRun() {
	i.muRunning.Lock() // run if not running already
	if !i.running {
		i.running = true
		go i.run()
	}
	i.muRunning.Unlock()
}

// Stop stops the loop. Should be called only after Start
func (i *Integrator) Stop() {
	i.chExit <- struct{}{}
	return
}

// blocking wait until time or stop
func delay(d time.Duration, stop <-chan struct{}) {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-stop:
		return
	case <-timer.C:
	}
}

// Should be called from Start only
func (i *Integrator) run() {
	defer func() {
		i.muRunning.Lock()
		i.running = false
		i.muRunning.Unlock()
	}()

	for {
		// maybe stop? non-blocking
		select {
		case <-i.chExit:
			return
		default:
		}

		// call Next with a Publish closure. If ErrEmpty, return. Retry otherwise
		err := i.nexter.Next(func(c *model.Customer) error { return i.publisher.Publish(c) })
		switch {
		case err == ErrEmpty:
			return
		case err != nil:
			i.Logger.Debug(fmt.Sprintf("Error: %s, wating to repeat", err))
			delay(i.DelayRetry, i.chExit)
		}
	}

}
