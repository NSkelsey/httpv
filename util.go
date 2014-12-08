package httpv

import (
	"encoding/hex"

	"github.com/conformal/btcec"
)

func FakeKey() (*btcec.PrivateKey, *btcec.PublicKey, error) {
	// Decode a hex-encoded private key.
	pkBytes, err := hex.DecodeString("22a47fa09a223f2aa079edf85a7c2d4f87" +
		"20ee63e502ee2869afab7de234b80c")
	if err != nil {
		return nil, nil, err
	}
	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
	return privKey, pubKey, nil
}
