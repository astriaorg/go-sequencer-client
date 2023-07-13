package client

import (
	"crypto/ed25519"
	"crypto/sha256"

	sqproto "github.com/astriaorg/go-sequencer-client/proto"
	proto "google.golang.org/protobuf/proto"
)

type Signer struct {
	private ed25519.PrivateKey
}

func GenerateSigner() (*Signer, error) {
	_, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}

	return &Signer{
		private: priv,
	}, nil
}

func (s *Signer) SignTransaction(tx *sqproto.UnsignedTransaction) (*sqproto.SignedTransaction, error) {
	bytes, err := proto.Marshal(tx)
	if err != nil {
		return nil, err
	}

	msg := sha256.Sum256(bytes)
	sig := ed25519.Sign(s.private, msg[:])
	return &sqproto.SignedTransaction{
		Transaction: tx,
		Signature:   sig,
		PublicKey:   s.private.Public().(ed25519.PublicKey),
	}, nil
}
