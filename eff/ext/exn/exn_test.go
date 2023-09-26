package exn_test

import (
	// "errors"
	"errors"
	"testing"

	"github.com/nymphium/eff.go/eff/ext/exn"

	"github.com/stretchr/testify/require"
)

func TestExn(t *testing.T) {
	t.Parallel()

	Err := errors.New("error")
	type ty = int

	t.Run("nothing", func(t *testing.T) {
		t.Parallel()

		res := exn.Try(func() ty {
			return 42
		}, exn.Handlers[ty]{
			Err: func(error) ty { return 0 },
		})

		require.Equal(t, 42, res)
	})

	t.Run("handle", func(t *testing.T) {
		t.Parallel()

		res := exn.Try(func() ty {
			return exn.Raise(Err).(ty)
		}, exn.Handlers[ty]{
			Err: func(error) ty { return 42 },
		})

		require.Equal(t, 42, res)

	})

	t.Run("joined error", func(t *testing.T) {
		t.Parallel()

		Err2 := errors.Join(Err, errors.New("error2"))

		res := exn.Try(func() ty {
			return exn.Raise(Err2).(ty)
		}, exn.Handlers[ty]{
			Err2: func(error) ty { return 42 },
		})

		require.Equal(t, 42, res)
	})

	t.Run("joined error", func(t *testing.T) {
		t.Parallel()

		Err2 := errors.Join(Err, errors.New("error2"))

		res := exn.Try(func() ty {
			return exn.Raise(Err2).(ty)
		}, exn.Handlers[ty]{
			Err: func(error) ty { return 42 },
		})

		require.Equal(t, 42, res)
	})

	t.Run("no reaching return", func(t *testing.T) {
		t.Parallel()

		res := exn.Try(func() ty {
			exn.Raise(Err)
			return 0
		}, exn.Handlers[ty]{
			Err: func(error) ty { return 42 },
		})

		require.Equal(t, 42, res)
	})

	t.Run("resend", func(t *testing.T) {
		t.Parallel()

		ErrX := errors.New("errorX")

		res := exn.Try(func() error {
			exn.Try(func() ty {
				return exn.Raise(Err).(ty)
			}, exn.Handlers[ty]{
				Err: func(error) ty {
					return exn.Raise(ErrX).(ty)
				},
			})

			return nil
		}, exn.Handlers[error]{
			ErrX: func(e error) error {
				return e
			},
		})

		require.ErrorIs(t, ErrX, res)
	})
}
