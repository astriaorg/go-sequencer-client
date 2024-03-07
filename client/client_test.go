package client

import (
	"context"
	"crypto/sha256"
	"testing"

	sqproto "buf.build/gen/go/astria/astria/protocolbuffers/go/astria/sequencer/v1alpha1"

	"github.com/stretchr/testify/require"
)

func TestSignAndBroadcastTx(t *testing.T) {
	signer, err := GenerateSigner()
	require.NoError(t, err)

	client, err := NewClient("http://localhost:26657")
	require.NoError(t, err)

	rollupId := sha256.Sum256([]byte("test-chain"))
	tx := &sqproto.UnsignedTransaction{
		Nonce: 1,
		Actions: []*sqproto.Action{
			{
				Value: &sqproto.Action_SequenceAction{
					SequenceAction: &sqproto.SequenceAction{
						RollupId: rollupId[:],
						Data:     []byte("test-data"),
					},
				},
			},
		},
	}

	signed, err := signer.SignTransaction(tx)
	require.NoError(t, err)

	resp, err := client.BroadcastTxSync(context.Background(), signed)
	require.NoError(t, err)
	require.Equal(t, resp.Code, uint32(0), resp.Log)
}

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
