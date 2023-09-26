package coro_test

import (
	"errors"
	"testing"

	"github.com/nymphium/eff.go/coro"
	"github.com/stretchr/testify/require"
)

func TestCoro(t *testing.T) {
	t.Parallel()

	t.Run("nothing", func(t *testing.T) {
		t.Parallel()

		c := coro.New(func(any) (any, error) {
			return 42, nil
		})

		res, err := c.Resume(nil)
		require.NoError(t, err)
		require.Equal(t, 42, res)
	})

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
		require.NoError(t, err)
		require.Equal(t, msg1, out)

		out, err = c.Resume(nil)
		require.NoError(t, err)
		require.Equal(t, msg2, out)

		out, err = c.Resume(nil)
		require.NoError(t, err)
		require.Equal(t, msg3, out)

		require.Equal(t, coro.StatusDone, c.Status)
	})

	t.Run("communicate", func(t *testing.T) {
		t.Parallel()

		syn := "syn"
		ack := "ack"

		c := coro.New(func(arg any) (any, error) {
			rcv := coro.Yield(syn)
			require.Equal(t, ack, rcv)

			return nil, nil
		})

		rcv, err := c.Resume(nil)
		require.NoError(t, err)
		require.Equal(t, syn, rcv)

		rcv, err = c.Resume(ack)
		require.NoError(t, err)
		require.Nil(t, rcv)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		e := errors.New("error inside")

		c := coro.New(func(any) (any, error) {
			return nil, e
		})

		_, err := c.Resume(nil)
		require.True(t, errors.Is(err, e))
	})

	t.Run("panic", func(t *testing.T) {
		t.Parallel()

		c := coro.New(func(any) (any, error) {
			panic(42)
		})

		require.Panics(t, func() { c.Resume(nil) })
	})

	t.Run("across call stack", func(t *testing.T) {
		t.Parallel()

		helloYield := func() any {
			return coro.Yield("hello")
		}

		c := coro.New(func(any) (v any, err error) {
			r := helloYield()
			require.Equal(t, "hello", r)
			return
		})

		_, err := c.Resume(nil)
		require.NoError(t, err)
	})

	t.Run("toplevel yield", func(t *testing.T) {
		t.Parallel()

		require.PanicsWithError(t, coro.ErrCoroutineCannotYield.Error(), func() { coro.Yield("hello") })
	})
}
