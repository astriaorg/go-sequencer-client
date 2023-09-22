package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	client "github.com/astriaorg/go-sequencer-client/client"
)

func Execute() {
	if len(os.Args) < 2 {
		fmt.Println("Expected a command.")
		printHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "createaccount":
		handleCreateAccount()
	case "getbalance":
		handleGetBalance()
	default:
		printHelp()
		os.Exit(1)
	}

}

func printHelp() {
	fmt.Println("Usage: go-sequencer-client-cli <command>")
	fmt.Println("Available commands are:")
	fmt.Println("  createaccount                       : creates an account")
	fmt.Println("  getbalance <rpc_endpoint> <address> : gets the balance of an account")
	os.Exit(1)
}

func handleCreateAccount() {
	privateKey, err := client.NewAccount()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// NOTE - creating a Signer to easily get public key and address
	signer := client.NewSigner(privateKey)
	address := signer.Address()
	fmt.Println("Created account:")
	fmt.Println("  Private Key:", hex.EncodeToString(privateKey.Seed()))
	fmt.Println("  Public Key: ", hex.EncodeToString(signer.PublicKey()))
	fmt.Println("  Address:    ", hex.EncodeToString(address[:]))
	os.Exit(0)
}

func handleGetBalance() {
	if len(os.Args) < 4 {
		fmt.Println("Expected an endpoint and address.")
		printHelp()
		os.Exit(1)
	}

	endpoint := os.Args[2]

	address, err := hex.DecodeString(os.Args[3])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client, err := client.NewClient(endpoint)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var address20 [20]byte
	copy(address20[:], address)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	balance, err := client.GetBalance(ctx, address20)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Balance:", balance)
	os.Exit(0)
}
