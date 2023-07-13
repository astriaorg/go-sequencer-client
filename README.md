# go-sequencer-client


### Generating Go files from protos

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
protoc -I=. --go_out=. types.proto transaction.proto && mv github.com/astriaorg/go-sequencer-client/* proto/ && rm -r github.com/
```
