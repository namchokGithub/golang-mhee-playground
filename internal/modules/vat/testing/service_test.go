package vat_testing

import (
	"proundmhee/internal/modules/vat"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestService_Calculate_DefaultRate(t *testing.T) {
	svc := vat.NewService()

	vat, total, err := svc.Calculate(1000, 0)
	require.NoError(t, err)
	require.Equal(t, 70.0, vat)
	require.Equal(t, 1070.0, total)
}

func TestService_Calculate_CustomRate(t *testing.T) {
	svc := vat.NewService()

	vat, total, err := svc.Calculate(200, 10)
	require.NoError(t, err)
	require.Equal(t, 20.0, vat)
	require.Equal(t, 220.0, total)
}
