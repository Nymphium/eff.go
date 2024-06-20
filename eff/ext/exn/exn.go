/*
exn provudes exception handling.

	  exn.Try(func() error {
			user := &User{Name:"Alice", Age: 20}
			res, err := json.Marshal(user)
			exn.Check(err)

			...

			return nil
		}, exn.Handlers[struct{}]{
			ErrAny: func(e error) (struct{}, error) {
				return struct{}, e
			}
		})
*/
package exn

import (
	"errors"
	"sort"

	"github.com/nymphium/eff.go/eff"
)

// The effect of exception.
var exn = eff.New()

// ErrAny is used for handling any error with lowest priority.
var ErrAny = errors.New("any error")

// The type of exception handlers.
type Handlers[Returning any] map[error]func(error) (Returning, error)

// Raise raises an exception.
func Raise(e error) any {
	return exn.Perform(e)
}

// Check checks an error.
// If the error is not nil, it raises the error.
func Check(e error) {
	if e != nil {
		Raise(e)
	}
}

/*
Try executes a function and catches an exception.
The handlers consider the joined errors.

	Err1 := errors.New("error1")
	Err2 := errors.Join(Err1, errors.New("error2"))

	exn.Try(func() ty {
	  return exn.Raise(Err2).(ty)
	}, exn.Handlers[ty]{
	  Err1: func(error) (ty, error) { return 42, nil },
	})
	// => 42, nil
*/
func Try[Returning any](f func() Returning, handlers Handlers[Returning]) (Returning, error) {
	t, err := eff.NewHandler(func(v any) (any, error) { return v, nil }).
		On(exn, func(e any, _ eff.Cont) (any, error) {
			err := e.(error)

			var anyhandler func(error) (Returning, error)
			errs := make(errs, 0, len(handlers))

			// handles joined error like exception and its class in Java
			for erk := range handlers {
				if errors.Is(erk, ErrAny) {
					anyhandler = handlers[erk]
				} else {
					errs = append(errs, erk)
				}
			}
			sort.Sort(&errs)

			for _, erk := range errs {
				if errors.Is(err, erk) {
					return handlers[erk](err)
				}
			}

			if anyhandler != nil {
				return anyhandler(err)
			}

			return Raise(err), nil
		}).
		Handle(func() (any, error) { return f(), nil })

	return t.(Returning), err
}
