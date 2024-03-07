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

// Should this live here?
type BalanceResponse struct {
	Denom   string   `json:"denom,omitempty"`
	Balance *big.Int `json:"balance,omitempty"`
}

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

func (c *Client) GetBalances(ctx context.Context, addr [20]byte) ([]*BalanceResponse, error) {
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

	protoBalanceResp := &sqproto.BalanceResponse{}
	err = proto.Unmarshal(resp.Response.Value, protoBalanceResp)
	if err != nil {
		return nil, err
	}

	return balanceResponseFromProto(protoBalanceResp), nil
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

func balanceResponseFromProto(resp *sqproto.BalanceResponse) []*BalanceResponse {
	var balanceResponses []*BalanceResponse
	for _, balance := range resp.Balances {
		balanceResponses = append(balanceResponses, &BalanceResponse{
			Balance: protoU128ToBigInt(balance.Balance),
			Denom:   balance.Denom,
		})
	}
	return balanceResponses
}
