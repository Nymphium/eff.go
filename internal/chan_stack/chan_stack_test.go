package chan_stack_test

import (
	"sync"
	"testing"
	"time"

	"github.com/nymphium/eff.go/internal/chan_stack"
	"github.com/stretchr/testify/require"
)

func TestChanStack(t *testing.T) {
	t.Parallel()

	t.Run("push and pop", func(t *testing.T) {
		t.Parallel()

		stack := chan_stack.New[int]()
		ch1 := make(chan int, 1)
		v1 := 42
		ch1 <- v1
		stack.Push(ch1)
		ch2 := make(chan int, 1)
		v2 := 43
		ch2 <- v2
		stack.Push(ch2)
		require.Equal(t, v2, <-stack.Pop())
		require.Equal(t, v1, <-stack.Pop())
		require.Nil(t, stack.Pop())
	})

	t.Run("delete by self", func(t *testing.T) {
		t.Parallel()

		stack := chan_stack.New[int]()
		ch1 := make(chan int, 1)
		v1 := 42
		ch1 <- v1
		stack.Push(ch1)
		ch2 := make(chan int, 1)
		v2 := 43
		ch2 <- v2
		stack.Push(ch2)
		stack.DeleteBySelf(ch1)
		require.Equal(t, v2, <-stack.Pop())
		require.Nil(t, stack.Pop())
	})

	t.Run("push and pop with goroutines", func(t *testing.T) {
		t.Parallel()

		stack := chan_stack.New[int]()
		ch1 := make(chan int, 1)
		v1 := 42
		ch1 <- v1
		ch2 := make(chan int, 1)
		v2 := 43
		ch2 <- v2

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			stack.Push(ch1)
		}()

		go func() {
			defer wg.Done()
			stack.Push(ch2)
		}()

		wg.Wait()

		require.Equal(t, v1+v2, <-stack.Pop()+<-stack.Pop())
		require.Nil(t, stack.Pop())
	})

	t.Run("delete by self with goroutines", func(t *testing.T) {
		t.Parallel()

		stack := chan_stack.New[int]()
		ch1 := make(chan int, 1)
		v1 := 42
		ch1 <- v1
		ch2 := make(chan int, 1)
		v2 := 43
		ch2 <- v2

		var wg sync.WaitGroup
		wg.Add(3)

		go func() {
			defer wg.Done()
			stack.Push(ch1)
		}()

		go func() {
			defer wg.Done()
			stack.Push(ch2)
		}()

		go func() {
			defer wg.Done()
			time.Sleep(1 * time.Second) // Ensure the channels have been pushed
			stack.DeleteBySelf(ch1)
		}()

		wg.Wait()

		require.Equal(t, v2, <-stack.Pop())
		require.Nil(t, stack.Pop())
	})
}
