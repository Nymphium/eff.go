package exn_test

import (
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

		res, err := exn.Try(func() ty {
			return 42
		}, exn.Handlers[ty]{
			Err: func(error) (ty, error) { return 0, nil },
		})

		require.NoError(t, err)
		require.Equal(t, 42, res)
	})

	t.Run("handle", func(t *testing.T) {
		t.Parallel()

		res, err := exn.Try(func() ty {
			return exn.Raise(Err).(ty)
		}, exn.Handlers[ty]{
			Err: func(error) (ty, error) { return 42, nil },
		})

		require.NoError(t, err)
		require.Equal(t, 42, res)

	})

	t.Run("joined error", func(t *testing.T) {
		t.Parallel()

		Err2 := errors.Join(Err, errors.New("error2"))

		res, err := exn.Try(func() ty {
			return exn.Raise(Err2).(ty)
		}, exn.Handlers[ty]{
			Err2: func(error) (ty, error) { return 42, nil },
		})

		require.NoError(t, err)
		require.Equal(t, 42, res)
	})

	t.Run("joined error", func(t *testing.T) {
		t.Parallel()

		Err2 := errors.Join(Err, errors.New("error2"))

		res, err := exn.Try(func() ty {
			return exn.Raise(Err2).(ty)
		}, exn.Handlers[ty]{
			Err: func(error) (ty, error) { return 42, nil },
		})

		require.NoError(t, err)
		require.Equal(t, 42, res)
	})

	t.Run("no reaching return", func(t *testing.T) {
		t.Parallel()

		res, err := exn.Try(func() ty {
			exn.Raise(Err)
			return 0
		}, exn.Handlers[ty]{
			Err: func(error) (ty, error) { return 42, nil },
		})

		require.NoError(t, err)
		require.Equal(t, 42, res)
	})

	t.Run("resend", func(t *testing.T) {
		t.Parallel()

		ErrX := errors.New("errorX")

		_, err := exn.Try(func() struct{} {
			exn.Try(func() ty {
				return exn.Raise(Err).(ty)
			}, exn.Handlers[ty]{
				Err: func(error) (ty, error) {
					return exn.Raise(ErrX).(ty), nil
				},
			})

			return struct{}{}
		}, exn.Handlers[struct{}]{
			ErrX: func(e error) (struct{}, error) {
				return struct{}{}, e
			},
		})

		require.ErrorIs(t, ErrX, err)
	})

	t.Run("anyerror", func(t *testing.T) {
		t.Parallel()

		res, err := exn.Try(func() ty {
			exn.Check(Err)
			return 0
		}, exn.Handlers[ty]{
			exn.ErrAny: func(error) (ty, error) { return 42, nil },
		})

		require.NoError(t, err)
		require.Equal(t, 42, res)
	})
}
