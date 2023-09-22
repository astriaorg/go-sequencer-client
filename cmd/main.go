package cmd

import (
	"encoding/hex"
	"fmt"
	"os"

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
	default:
		printHelp()
		os.Exit(1)
	}

}

func printHelp() {
	fmt.Println("Usage: go-sequencer-client-cli <command>")
	fmt.Println("Available commands are:")
	fmt.Println("  createaccount: creates an account")
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
