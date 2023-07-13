package client

import (
	"context"
	"testing"

	sqproto "github.com/astriaorg/go-sequencer-client/proto"

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
						ChainId: []byte("test-chain"),
						Data:    []byte("test-data"),
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
