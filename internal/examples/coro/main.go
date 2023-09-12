package main

import (
	"fmt"

	"github.com/nymphium/eff.go/coro"
)

func f(label string, arg, ret any) any {
	fmt.Println(label, arg)
	return coro.Yield(ret)
}

func main() {
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
}
