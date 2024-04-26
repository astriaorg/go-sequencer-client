package client

import (
	"context"
	"crypto/sha256"
	"testing"

	primproto "buf.build/gen/go/astria/primitives/protocolbuffers/go/astria/primitive/v1"
	txproto "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria/protocol/transactions/v1alpha1"

	"github.com/stretchr/testify/require"
)

func TestSignAndBroadcastTx(t *testing.T) {
	signer, err := GenerateSigner()
	require.NoError(t, err)

	client, err := NewClient("http://localdev.me:26657")
	require.NoError(t, err)

	rollupIdBytes := sha256.Sum256([]byte("test-rollup"))
	rollupId := &primproto.RollupId{
		Inner: rollupIdBytes[:],
	}
	tx := &txproto.UnsignedTransaction{
		Params: &txproto.TransactionParams{
			ChainId: "test-chain",
			Nonce:   1,
		},
		Actions: []*txproto.Action{
			{
				Value: &txproto.Action_SequenceAction{
					SequenceAction: &txproto.SequenceAction{
						RollupId:   rollupId,
						Data:       []byte("test-data"),
						FeeAssetId: DefaultAstriaAssetID[:],
					},
				},
			},
		},
	}

	signed, err := signer.SignTransaction(tx)
	require.NoError(t, err)

	resp, err := client.BroadcastTxSync(context.Background(), signed)
	require.NoError(t, err)
	require.Equal(t, uint32(0), resp.Code, resp.Log)
}

func TestGetBalance(t *testing.T) {
	client, err := NewClient("http://localdev.me:26657")
	require.NoError(t, err)

	balance, err := client.GetBalances(context.Background(), [20]byte{})
	require.NoError(t, err)
	require.Empty(t, balance)
}

func TestGetNonce(t *testing.T) {
	client, err := NewClient("http://localdev.me:26657")
	require.NoError(t, err)

	nonce, err := client.GetNonce(context.Background(), [20]byte{})
	require.NoError(t, err)
	require.Equal(t, nonce, uint32(0))
}
