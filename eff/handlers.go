package eff

import "github.com/nymphium/eff.go/coro"

type handleFunc = func(arg any, k Cont) (any, error)

type handler struct {
	effh map[effID]handleFunc
	valh Cont
}

// NewHandler creates a new handler with a default value handler.
func NewHandler(valh Cont) handler {
	return handler{
		effh: map[effID]handleFunc{},
		valh: valh,
	}
}

// On chains a handler for the effect.
func (h handler) On(eff *T, handleFunc handleFunc) handler {
	h.effh[eff.id] = handleFunc
	return h
}

// To updates a value handler to the handler.
func (h handler) To(valh Cont) handler {
	h.valh = valh
	return h
}

func (h handler) Handle(thunk func() (any, error)) (any, error) {
	if h.valh == nil {
		h.valh = func(r any) (any, error) {
			return r, nil
		}
	}

	co := coro.New(func(any) (any, error) {
		return thunk()
	})

	var (
		cont     Cont
		handle   Cont
		rehandle func(Cont) Cont
	)

	cont = func(a any) (r any, err error) {
		rr, err := co.Resume(a)
		if err != nil {
			return r, err
		}
		return handle(rr)
	}

	rehandle = func(k Cont) Cont {
		hnew := h
		return func(a any) (any, error) {
			return hnew.To(cont).
				Handle(func() (any, error) {
					return k(a)
				})
		}
	}

	handle = func(r any) (any, error) {
		switch r := r.(type) {
		case eff:
			effh, ok := h.effh[r.id]
			if ok {
				return effh(r.arg, cont)
			} else {
				return r.resend(cont), nil
			}
		case resend:
			effh, ok := h.effh[r.eff.id]
			if ok {
				return effh(r.arg, rehandle(r.cont))
			} else {
				return r.resend(rehandle(r.cont)), nil
			}
		default:
			return h.valh(r)
		}
	}

	return cont(nil)
}
