package coro_test

import (
	"errors"
	"testing"

	"github.com/nymphium/eff/src/coro"
	"github.com/stretchr/testify/assert"
)

func TestCoro(t *testing.T) {
	t.Parallel()

	t.Run("basic", func(t *testing.T) {
		t.Parallel()

		msg1 := "hello"
		msg2 := "world"
		msg3 := "!"

		c := coro.New(func(_ any, yield coro.Yield[any, string]) (out string, err error) {
			_, _ = yield(msg1)
			_, _ = yield(msg2)

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

		assert.Equal(t, coro.StateDone, c.State())
	})

	t.Run("communicate", func(t *testing.T) {
		t.Parallel()

		syn := "syn"
		ack := "ack"

		c := coro.New(func(_ any, yield coro.Yield[any, any]) (out any, err error) {
			rcv, err := yield(syn)
			assert.NoError(t, err)
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

	t.Run("cancel", func(t *testing.T) {
		t.Parallel()

		c := coro.New(func(_ any, yield coro.Yield[any, any]) (out any, err error) {
			_, err = yield(nil)
			assert.True(t, errors.Is(err, coro.ErrCanceled), "got", err)

			return nil, nil
		})

		_, _ = c.Resume(nil)
		err := c.Cancel()
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		e := errors.New("error inside")

		c := coro.New(func(_ any, yield coro.Yield[any, any]) (any, error) {
			return nil, e
		})

		_, err := c.Resume(nil)
		assert.True(t, errors.Is(err, e))
	})

	t.Run("panic", func(t *testing.T) {
		t.Parallel()

		c := coro.New(func(_ any, yield coro.Yield[any, any]) (any, error) {
			panic(42)
		})

		assert.Panics(t, func() { c.Resume(nil) })
	})
}
