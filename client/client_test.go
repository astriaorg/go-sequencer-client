package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBalance(t *testing.T) {
	client, err := NewClient("http://localhost:26657")
	require.NoError(t, err)

	balance, err := client.GetBalances(context.Background(), [20]byte{})
	require.NoError(t, err)
	require.Empty(t, balance)
}

func TestGetNonce(t *testing.T) {
	client, err := NewClient("http://localhost:26657")
	require.NoError(t, err)

	nonce, err := client.GetNonce(context.Background(), [20]byte{})
	require.NoError(t, err)
	require.Equal(t, nonce, uint32(0))
}
