/*
One-shot algebraic effect using coroutines.

  defering := eff.New()

  h := eff.NewHandler(defering, func(r any) (any, error) {
    print("!")
    return r, nil
  }}.
    On(defering, func(th any, k eff.Cont) (any, error) {
      k(nil)
      th()
      return nil
    })

  h.Handle(func() (any, error) {
    defering.Perform(func() {
      print("world")
    })
    print("hello, ")
    return nil, nil
  })
*/
package eff

import "github.com/nymphium/eff.go/coro"

// The (rough) type of continuation.
type Cont = func(any) (any, error)

// The ID of effect.
// Each effect is unique.
type effID *struct{}

// The type of effect.
type T struct{ id effID }

type eff struct {
	*T
	arg any
}

type resend struct {
	eff
	cont func(any) (any, error)
}

// Resend the effect invocation with its continuation to the next outer handler of the "current" one.
func (e eff) resend(k func(any) (any, error)) any {
	return coro.Yield(resend{e, k})
}

// Perform the effect with an argument.
// The corresponding handler will catch the arguemnt.
func (t *T) Perform(arg any) any {
	eff := eff{t, arg}
	return coro.Yield(eff)
}

// Create a new effect.
func New() *T {
	return &T{id: &struct{}{}}
}
