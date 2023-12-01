package client

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/cometbft/cometbft/rpc/client"
	"github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	proto "google.golang.org/protobuf/proto"

	primproto "buf.build/gen/go/astria/astria/protocolbuffers/go/astria/primitive/v1"
	sqproto "buf.build/gen/go/astria/astria/protocolbuffers/go/astria/sequencer/v1alpha1"
)

// Client is an HTTP tendermint client.
type Client struct {
	client *http.HTTP
}

func NewClient(url string) (*Client, error) {
	client, err := http.New(url, "")
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

func (c *Client) GetBalance(ctx context.Context, addr [20]byte) (*big.Int, error) {
	resp, err := c.client.ABCIQueryWithOptions(ctx, fmt.Sprintf("accounts/balance/%x", addr), []byte{}, client.ABCIQueryOptions{
		Height: 0,
		Prove:  false,
	})
	if err != nil {
		return nil, err
	}

	if resp.Response.Code != 0 {
		return nil, errors.New(resp.Response.Log)
	}

	balanceResp := &sqproto.BalanceResponse{}
	err = proto.Unmarshal(resp.Response.Value, balanceResp)
	if err != nil {
		return nil, err
	}

	return protoU128ToBigInt(balanceResp.Balance), nil
}

func (c *Client) GetNonce(ctx context.Context, addr [20]byte) (uint32, error) {
	resp, err := c.client.ABCIQueryWithOptions(ctx, fmt.Sprintf("accounts/nonce/%x", addr), []byte{}, client.ABCIQueryOptions{
		Height: 0,
		Prove:  false,
	})
	if err != nil {
		return 0, err
	}

	if resp.Response.Code != 0 {
		return 0, errors.New(resp.Response.Log)
	}

	nonceResp := &sqproto.NonceResponse{}
	err = proto.Unmarshal(resp.Response.Value, nonceResp)
	if err != nil {
		return 0, err
	}

	return nonceResp.Nonce, nil
}

func protoU128ToBigInt(u128 *primproto.Uint128) *big.Int {
	lo := big.NewInt(0).SetUint64(u128.Lo)
	hi := big.NewInt(0).SetUint64(u128.Hi)
	hi.Lsh(hi, 64)
	return lo.Add(lo, hi)
}
