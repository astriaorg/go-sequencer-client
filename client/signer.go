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

func NewSigner(private ed25519.PrivateKey) *Signer {
	return &Signer{
		private: private,
	}
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

	sig := ed25519.Sign(s.private, bytes)
	return &sqproto.SignedTransaction{
		Transaction: tx,
		Signature:   sig,
		PublicKey:   s.private.Public().(ed25519.PublicKey),
	}, nil
}

// Seed returns the 32-byte "seed" for the key, which is used as the
// input to generate a private key in the rust implementation, ie:
// `ed25519_consensus::SigningKey::from(seed)`
func (s *Signer) Seed() [ed25519.SeedSize]byte {
	return [ed25519.SeedSize]byte(s.private.Seed())
}

func (s *Signer) PublicKey() ed25519.PublicKey {
	return s.private.Public().(ed25519.PublicKey)
}

func (s *Signer) Address() [20]byte {
	hash := sha256.Sum256(s.PublicKey())
	var addr [20]byte
	copy(addr[:], hash[:20])
	return addr
}
