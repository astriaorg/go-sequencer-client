package client

import (
	"context"

	"github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	proto "google.golang.org/protobuf/proto"

	sqproto "github.com/astriaorg/go-sequencer-client/proto"
)

// Client is an HTTP tendermint client.
type Client struct {
	client *http.HTTP
}

func NewClient(url string) (*Client, error) {
	client, err := http.New(url, "/websocket")
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) BroadcastTxSync(ctx context.Context, tx *sqproto.SignedTransaction) (*coretypes.ResultBroadcastTx, error) {
	bytes, err := proto.Marshal(tx)
	if err != nil {
		return nil, err
	}

	return c.client.BroadcastTxSync(ctx, bytes)
}
