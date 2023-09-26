package exn

import (
	"errors"
	"sort"

	"github.com/nymphium/eff.go/eff"
)

var Exn = eff.New()

type Handlers[Returning any] map[error]func(error) Returning

type errs []error

func (errs *errs) Len() int { return len(*errs) }
func (errs *errs) Swap(i, j int) {
	(*errs)[i], (*errs)[j] = (*errs)[j], (*errs)[i]
}

func (errs *errs) Less(i, j int) bool {
	e1 := (*errs)[i]
	e2 := (*errs)[j]
	e1Is := errors.Is(e1, e2)
	e2Is := errors.Is(e2, e1)
	switch {
	case e1Is && !e2Is:
		return false
	case !e1Is && e2Is:
		return true
	default:
		return false
	}
}

// Raise raises an exception.
func Raise(e error) any {
	return Exn.Perform(e)
}

// Try executes a function and catches an exception.
func Try[Returning any](f func() Returning, handlers Handlers[Returning]) Returning {
	// error should be nil
	t, _ := eff.NewHandler(func(v any) (any, error) { return v, nil }).
		On(Exn, func(e any, _ eff.Cont) (any, error) {
			err := e.(error)

			errs := make(errs, 0, len(handlers))

			// handles joined error like exception and its class in Java
			for erk := range handlers {
				errs = append(errs, erk)
			}
			sort.Sort(&errs)

			for _, erk := range errs {
				if errors.Is(err, erk) {
					return handlers[erk](err), nil
				}
			}

			return Raise(err), nil
		}).
		Handle(func() (any, error) { return f(), nil })

	return t.(Returning)
}
