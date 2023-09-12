// coro provides asytmmetric stackful coroutine
package coro

import (
	"errors"

	"github.com/nymphium/eff.go/internal/chan_stack"
)

type Status string

const (
	StatusNotStarted Status = "not started"
	StatusRunning    Status = "pending"
	StatusSuspending Status = "suspending"
	StatusDone       Status = "done"
)

var (
	ErrCoroutine             = errors.New("coroutine error")
	ErrCoroutineNotRunning   = errors.Join(ErrCoroutine, errors.New("coroutine not running"))
	ErrCoroutineCannotResume = errors.Join(ErrCoroutine, errors.New("coroutine cannot resume"))
)

// result of coroutine
type result struct {
	panic any
	error error
	value any
}

// coroutine object
type T struct {
	initialEntering chan any    // for initial argument of wrapped function
	returning       chan result // for final return value of wrapped function
	Status          Status
}

func New(f func(any) (any, error)) *T {
	initial := make(chan any)
	returning := make(chan result)

	co := &T{
		Status:          StatusNotStarted,
		initialEntering: initial,
		returning:       returning,
	}

	// push goroutine and wait for initial value
	go func() {
		m := result{}

		defer func() {
			if r := recover(); r != nil {
				m.panic = r
			}

			returning <- m
		}()

		m.value, m.error = f(<-initial)
	}()

	return co
}

func Yield(x any) any {
	// get returning point of resume
	returning := chan_stack.Pop()
	defer close(returning)

	// make entering point and push for next resume
	entering := chan_stack.PushNew()
	defer chan_stack.DeleteBySelf(entering)

	// pass value to current resume point
	returning <- x

	// suspend the coroutine and wait for next resume
	return <-entering
}

func (co *T) getEnteringPoint() (chan any, error) {
	switch co.Status {
	case StatusNotStarted:
		return co.initialEntering, nil
	case StatusSuspending:
		return chan_stack.Pop(), nil
	default:
		return nil, errors.Join(ErrCoroutineCannotResume, errors.New(string(co.Status)))
	}
}

func (co *T) raceReturningPoint(returning chan any) (any, error) {
	select {
	case v := <-returning:
		co.Status = StatusSuspending
		return v, nil
	case m := <-co.returning:
		if m.panic != nil {
			panic(m.panic)
		}

		co.Status = StatusDone
		return m.value, m.error
	}
}

func (co *T) Resume(x any) (any, error) {
	entering, err := co.getEnteringPoint()
	if err != nil {
		return nil, err
	}

	// make returning point and push for next yield
	returning := chan_stack.PushNew()
	defer chan_stack.DeleteBySelf(returning)

	co.Status = StatusRunning

	// pass value to suspending point
	entering <- x

	return co.raceReturningPoint(returning)
}
