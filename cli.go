package httpv

import (
	"encoding/base64"
	"errors"

	"github.com/NSkelsey/net/http"
	"github.com/conformal/btcec"
)

var (
	ErrRespNoVerHead = errors.New("httpv: Response does not have a Httpv-Ver Header")
	ErrRespNoTs      = errors.New("httpv: Response does not have a Httpv-Ts Header")
	ErrRespNoSig     = errors.New("httpv: Response does not have a Httpv-Sig Header")
	ErrRespBadNonce  = errors.New("httpv: Response has a bad nonce")
)

// Assert that resp has the correct headers
func assertRespHeader(header http.Header) error {
	if header.Get("Httpv-Ver") != "0.1" {
		return ErrRespNoVerHead
	}
	if header.Get("Httpv-Ts") == "" {
		return ErrRespNoTs
	}
	if header.Get("Httpv-Sig") == "" {
		return ErrRespNoSig
	}
	if header.Get("Httpv-N") == "" {
		return ErrRespBadNonce
	}
	return nil
}

// Pull out sig from the http headers
func getSig(resp *http.Response) (*btcec.Signature, error) {
	sig_str := resp.Header.Get("Httpv-Sig")
	sigBytes, err := base64.StdEncoding.DecodeString(sig_str)
	if err != nil {
		return nil, err
	}

	sig, err := btcec.ParseSignature(sigBytes, btcec.S256())
	if err != nil {
		return nil, err
	}

	return sig, nil
}
