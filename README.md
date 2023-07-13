# go-sequencer-client

### Usage
```go
package main

import (
	"context"
	"fmt"

	client "github.com/astriaorg/go-sequencer-client"
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

### Generating Go files from protos

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
protoc -I=. --go_out=. types.proto transaction.proto && mv github.com/astriaorg/go-sequencer-client/* proto/ && rm -r github.com/
```
