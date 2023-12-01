package client

import (
	"context"
	"math/big"
	"testing"

	sqproto "buf.build/gen/go/astria/astria/protocolbuffers/go/astria/sequencer/v1alpha1"

	"github.com/stretchr/testify/require"
)

func TestSignAndBroadcastTx(t *testing.T) {
	signer, err := GenerateSigner()
	require.NoError(t, err)

	client, err := NewClient("http://localhost:26657")
	require.NoError(t, err)

	tx := &sqproto.UnsignedTransaction{
		Nonce: 1,
		Actions: []*sqproto.Action{
			{
				Value: &sqproto.Action_SequenceAction{
					SequenceAction: &sqproto.SequenceAction{
						RollupId: []byte("test-chain"),
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

	balance, err := client.GetBalance(context.Background(), [20]byte{})
	require.NoError(t, err)
	require.Equal(t, balance, big.NewInt(0))
}

func TestGetNonce(t *testing.T) {
	client, err := NewClient("http://localhost:26657")
	require.NoError(t, err)

	nonce, err := client.GetNonce(context.Background(), [20]byte{})
	require.NoError(t, err)
	require.Equal(t, nonce, uint32(0))
}
