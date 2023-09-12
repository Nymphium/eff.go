package coro_test

import (
	"errors"
	"testing"

	"github.com/nymphium/eff.go/coro"
	"github.com/stretchr/testify/assert"
)

func TestCoro(t *testing.T) {
	t.Parallel()

	t.Run("basic", func(t *testing.T) {
		t.Parallel()

		msg1 := "hello"
		msg2 := "world"
		msg3 := "!"

		c := coro.New(func(arg any) (out any, err error) {
			_ = coro.Yield(msg1)
			_ = coro.Yield(msg2)

			return msg3, nil
		})

		out, err := c.Resume(nil)
		assert.NoError(t, err)
		assert.Equal(t, msg1, out)

		out, err = c.Resume(nil)
		assert.NoError(t, err)
		assert.Equal(t, msg2, out)

		out, err = c.Resume(nil)
		assert.NoError(t, err)
		assert.Equal(t, msg3, out)

		assert.Equal(t, coro.StatusDone, c.Status)
	})

	t.Run("communicate", func(t *testing.T) {
		t.Parallel()

		syn := "syn"
		ack := "ack"

		c := coro.New(func(arg any) (any, error) {
			rcv := coro.Yield(syn)
			assert.Equal(t, ack, rcv)

			return nil, nil
		})

		rcv, err := c.Resume(nil)
		assert.NoError(t, err)
		assert.Equal(t, syn, rcv)

		rcv, err = c.Resume(ack)
		assert.NoError(t, err)
		assert.Nil(t, rcv)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		e := errors.New("error inside")

		c := coro.New(func(any) (any, error) {
			return nil, e
		})

		_, err := c.Resume(nil)
		assert.True(t, errors.Is(err, e))
	})

	t.Run("panic", func(t *testing.T) {
		t.Parallel()

		c := coro.New(func(any) (any, error) {
			panic(42)
		})

		assert.Panics(t, func() { c.Resume(nil) })
	})

	t.Run("across call stack", func(t *testing.T) {
		t.Parallel()

		helloYield := func() any {
			return coro.Yield("hello")
		}

		c := coro.New(func(any) (v any, err error) {
			r := helloYield()
			assert.Equal(t, "hello", r)
			return
		})

		_, err := c.Resume(nil)
		assert.NoError(t, err)
	})
}
