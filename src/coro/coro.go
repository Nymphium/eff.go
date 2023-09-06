// This package is based on Russ Cox's coroutines
// ref: https://research.swtch.com/coro
package coro

import (
	"fmt"

	"errors"
)

// the state of the coroutine.
// - StateNotRunning is the initial state of a coroutine.
// - StateRunning is the state of a coroutine while it is running.
// - StateSuspending is the state of a coroutine while it is suspended.
// - StateCanceled is the state of a coroutine while it is canceled.
// - StateCanceled is the state of a coroutine while it is panicked.
// - StateDone is the state of a coroutine after it has finished running.
type State string

const (
	StateNotRunning State = "not running"
	StateRunning    State = "running"
	StateSuspending State = "suspending"
	StateCanceled   State = "canceled"
	StatePanicked   State = "panicked"
	StateDone       State = "done"
)

var (
	Err             = errors.New("coroutine error")
	ErrCanceled     = errors.Join(Err, errors.New("coroutine canceled"))
	ErrCannotResume = errors.Join(Err, errors.New("coroutine cannot resume"))
)

type T[In, Out any] interface {
	// Resume runs the coroutine with passing input to it until it yields or returns.
	// If the coroutine yields, then resume returns the value passed to Yield.
	// It may returns ErrCannotResume if the coroutine is not in the right state to resume.
	Resume(input In) (Out, error)

	//  Cancel stops the execution of f and shuts down the coroutine.
	// If resume has not been called, then f does not run at all.
	// Otherwise, cancel causes the blocked Yield call to panic with an error satisfying errors.Is(err, ErrCanceled).
	Cancel() error

	// See State type
	State() State
}

// t controls the execution of a coroutine.
type t[In, Out any] struct {
	cin   chan msg[In]
	cout  chan msg[Out]
	state State
}

// msg is the message type used to communicate between the coroutine and its caller.
type msg[T any] struct {
	panic any
	error error
	val   T
}

// Yield is a type of `yield`, which is used to suspend the coroutine and return a value `out` to the caller.
// `in` is the input value passed to the next `resume`.
// It may rerturn ErrCanceled if the coroutine is canceled.
type Yield[In, Out any] func(out Out) (in In, err error)

type Cell[In, Out any] func(In, Yield[In, Out]) (Out, error)

func (t *t[In, Out]) Resume(in In) (out Out, err error) {
	if t.state != StateNotRunning && t.state != StateSuspending {
		err = ErrCannotResume
		return
	}

	// enter resume point
	t.state = StateRunning
	t.cin <- msg[In]{val: in}

	m := <-t.cout
	if m.panic != nil {
		t.state = StatePanicked
		panic(m.panic)
	}

	return m.val, m.error
}

func (t *t[In, Out]) yield(out Out) (In, error) {

	// suspend coroutine here
	t.state = StateSuspending
	t.cout <- msg[Out]{val: out}

	// resume point
	m := <-t.cin
	if m.panic != nil {
		t.state = StatePanicked
		panic(m.panic)
	}

	return m.val, m.error
}

func (t *t[In, Out]) Cancel() (err error) {
	e := fmt.Errorf("%w", ErrCanceled)
	t.state = StateCanceled
	t.cin <- msg[In]{error: e}

	m := <-t.cout
	if m.panic != nil {
		t.state = StatePanicked
		panic(m.panic)
	}

	return m.error
}

// run wraps f and runs it in a separate goroutine.
func (t *t[In, Out]) run(f Cell[In, Out]) {
	// trap panic in the coroutine and send it back to the caller
	defer func() {
		if t.state == StateRunning {
			if r := recover(); r != nil {

				t.state = StatePanicked
				t.cout <- msg[Out]{panic: r}
			}
		}
	}()

	var (
		out Out
		err error
	)

	m := <-t.cin
	if m.panic == nil {
		out, err = f(m.val, t.yield)
	}

	t.state = StateDone
	t.cout <- msg[Out]{val: out, error: err}
}

func (t *t[In, Out]) State() State {
	return t.state
}

// creates a new, paused coroutine ready to run the function f.
// The new coroutine never runs on its own: it only runs by calling resume or cancel.
func New[In, Out any](f Cell[In, Out]) T[In, Out] {
	t := &t[In, Out]{
		cin:   make(chan msg[In]),
		cout:  make(chan msg[Out]),
		state: StateNotRunning,
	}

	go t.run(f)

	return t
}
