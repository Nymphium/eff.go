package exn

import "errors"

type errs []error

func (errs *errs) Len() int { return len(*errs) }
func (errs *errs) Swap(i, j int) {
	(*errs)[i], (*errs)[j] = (*errs)[j], (*errs)[i]
}
func (errs *errs) Less(i, j int) bool {
	e1, e2 := (*errs)[i], (*errs)[j]
	e1Is, e2Is := errors.Is(e1, e2), errors.Is(e2, e1)
	switch {
	case e1Is && !e2Is:
		return false
	case !e1Is && e2Is:
		return true
	default:
		return false
	}
}
