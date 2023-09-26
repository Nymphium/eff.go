package eff_test

import (
	"testing"

	"github.com/nymphium/eff.go/eff"
	"github.com/stretchr/testify/require"
)

func TestEff(t *testing.T) {
	t.Parallel()

	Eff := eff.New()
	h := eff.NewHandler(func(arg any) (out any, err error) {
		return arg, nil
	})

	t.Run("abort", func(t *testing.T) {
		t.Parallel()

		h := h.On(Eff, func(arg any, k eff.Cont) (any, error) {
			return arg, nil
		})

		res, err := h.Handle(func() (any, error) {
			return Eff.Perform(42), nil
		})

		require.NoError(t, err)
		require.Equal(t, 42, res)
	})

	t.Run("nested", func(t *testing.T) {
		t.Parallel()

		hin := h.On(Eff, func(arg any, k eff.Cont) (any, error) {
			return k(Eff.Perform(arg))
		})

		hout := h.On(Eff, func(arg any, k eff.Cont) (any, error) {
			return k(arg)
		})

		res, err := hout.Handle(func() (any, error) {
			return hin.Handle(func() (any, error) {
				return Eff.Perform(42), nil
			})
		})

		require.NoError(t, err)
		require.Equal(t, 42, res)
	})

	t.Run("nested abort", func(t *testing.T) {
		t.Parallel()

		hin := h.On(Eff, func(arg any, k eff.Cont) (any, error) {
			return Eff.Perform(arg), nil
		})

		hout := h.On(Eff, func(arg any, k eff.Cont) (any, error) {
			return arg, nil
		})

		res, err := hout.Handle(func() (any, error) {
			return hin.Handle(func() (any, error) {
				return Eff.Perform(42), nil
			})
		})

		require.NoError(t, err)
		require.Equal(t, 42, res)
	})
}
