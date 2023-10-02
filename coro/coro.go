/*
coro provides asytmmetric stackful coroutine

  var (
    co1 *coro.T
    co2 *coro.T
  )

  co1 = coro.New(func(arg any) (any, error) {
    fmt.Println("co1", arg)
    r, _ := co2.Resume("go co2")
    fmt.Println("co1", r)

    return nil, nil
  })

  co2 = coro.New(func(arg any) (any, error) {
    fmt.Println("co2", arg)
    return coro.Yield("go co1"), nil
  })

  co1.Resume("start")

  co3 := coro.New(func(arg any) (any, error) {
    r := f("co1", arg, 42)
    r = f("co2", r, 43)
    fmt.Println("co3", r)
    return 50, nil
  })

  fmt.Println("-----")

  r, _ := co3.Resume("start")
  fmt.Println("main1", r)
  r, _ = co3.Resume("res1")
  fmt.Println("main2", r)
  r, _ = co3.Resume("res2")
  fmt.Println("main3", r)
*/
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
	ErrCoroutineCannotYield  = errors.Join(ErrCoroutine, errors.New("coroutine cannot yield"))

	stack = chan_stack.New[result]()
)

// result of coroutine
type result struct {
	panic any
	error error
	value any
}

// coroutine object
type T struct {
	ch     chan result
	Status Status
}

func (co *T) validateResumableStatus() error {
	switch co.Status {
	case StatusNotStarted, StatusSuspending:
		return nil
	default:
		return errors.Join(ErrCoroutineCannotResume, errors.New(string(co.Status)))
	}
}

func New(f func(any) (any, error)) *T {
	ch := make(chan result)

	co := &T{
		Status: StatusNotStarted,
		ch:     ch,
	}

	// run goroutine and wait for initial value
	go func() {
		m := result{}

		defer func() {
			co.Status = StatusDone

			if r := recover(); r != nil {
				m.panic = r
			}

			ch <- m
		}()

		m.value, m.error = f(<-ch)
	}()

	return co
}

// Yield gets out from a coroutine with passing value to the caller
func Yield(x any) any {
	// get ch point of resume
	ch := stack.Pop()
	// toplevel yield not allowed
	if ch == nil {
		panic(ErrCoroutineCannotYield)
	}

	// pass value to current resume point
	ch <- result{value: x}

	// suspend the coroutine and wait for next resume
	return (<-ch).value
}

func (co *T) Resume(x any) (any, error) {
	if err := co.validateResumableStatus(); err != nil {
		return nil, err
	}

	ch := co.ch

	// make returning point and push for next yield
	stack.Push(ch)
	defer stack.DeleteBySelf(ch)

	// pass value to suspending point with status running
	co.Status = StatusRunning
	co.ch <- result{value: x}
	if co.Status == StatusRunning {
		// change status suspending if still running
		co.Status = StatusSuspending
	}

	res := <-ch
	if res.panic != nil {
		panic(res.panic)
	}
	return res.value, res.error
}
