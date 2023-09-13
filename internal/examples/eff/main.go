package main

import (
	"fmt"

	"github.com/nymphium/eff.go/eff"
)

func main() {
	greet := eff.New()

	h := eff.NewHandler(func(r any) (any, error) {
		fmt.Println("done")
		return r, nil
	}).
		On(greet, func(arg any, k eff.Cont) (any, error) {
			fmt.Println(arg)
			k("hi")
			fmt.Println("bye")
			return nil, nil
		})

	h.Handle(func() (any, error) {
		res := greet.Perform("hello")
		fmt.Println(">>", res)
		return "", nil
	})
}
