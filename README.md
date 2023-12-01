# go-sequencer-client

### Usage
```go
package main

import (
	"context"
	"fmt"

	client "github.com/astriaorg/go-sequencer-client/client"
	sqproto "github.com/astriaorg/go-sequencer-client/proto"
)

func main() {
	signer, err := client.GenerateSigner()
	if err != nil {
		panic(err)
	}

	// default tendermint RPC endpoint
	c, err := client.NewClient("http://localhost:26657")
	if err != nil {
		panic(err)
	}

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
	if err != nil {
		panic(err)
	}

	resp, err := c.BroadcastTxSync(context.Background(), signed)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
```

### CLI

```sh
# build the binary
go build -o go-sequencer-client-cli ./cmd/main.go
# run the binary
./go-sequencer-client-cli createaccount
./go-sequencer-client-cli getbalance 0f7681a7628cd931e0a036633168a96bb1a5f012 tcp://sequencer.localdev.me:80
```
