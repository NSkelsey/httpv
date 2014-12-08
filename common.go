package httpv

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/NSkelsey/net/http"
	"github.com/conformal/btcec"
)

var (
	HttpvVer string = "0.1"
)

type Conversation struct {
	servername string
	req        *http.Request
	resp       *http.Response
	privkey    *btcec.PrivateKey
	pubkey     *btcec.PublicKey
	hash       []byte
	debug      bool
}

func NewConversation(host string, pubkey *btcec.PublicKey, privkey *btcec.PrivateKey) Conversation {

	c := Conversation{
		servername: host,
		pubkey:     pubkey,
		privkey:    privkey,
		debug:      false,
	}
	return c
}

// Adds the httpv request into the conversation. If the request does not have the
// proper headers then an error is thrown.
func (c *Conversation) AddRequest(req http.Request) error {
	if c.req != nil {
		return errors.New("httpv: A Request was already added!")
	}

	if err := enforceReq(req); err != nil {
		return err
	}

	cleanHeader := reqHeaders()

	c.req = &http.Request{
		Method:     "GET",
		URL:        req.URL,
		Host:       c.servername,
		Proto:      req.Proto,
		ProtoMajor: req.ProtoMajor,
		ProtoMinor: req.ProtoMinor,
		Header:     cleanHeader,
	}
	if c.debug {
		fmt.Printf("========\n%s\n=======\n", c.req)
	}
	return nil
}

// Adds an http.Response to the conversation. The body and a select number of
// headers are copied over to the request that is actually signed.
func (c *Conversation) AddResponse(resp http.Response) error {
	if c.resp != nil {
		return errors.New("httpv: Conversation already has an initialized response!")
	}

	var header http.Header
	var err error
	if resp.Header.Get("Httpv-Sig") == "" {
		// The server is adding a new response that needs headers added
		header, err = respHeaders(resp.Header)
		if err != nil {
			return err
		}
	} else {
		// This must be a client side response, assert that the headers are
		// properly formed
		if err = assertRespHeader(resp.Header); err != nil {
			return err
		}
		header = resp.Header
	}

	c.resp = &http.Response{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Proto:      resp.Proto,
		ProtoMajor: resp.ProtoMajor,
		ProtoMinor: resp.ProtoMinor,
		Header:     header,
		Body:       resp.Body,
	}
	return nil
}

// Returns the signed http.Response to be sent back to a client. The headers
// that are actually signed is a proper subset of the ones sent over the wire.
// See respHeaders for the headers that actually get signed.
func (c *Conversation) EmitResponse() (*http.Response, error) {
	if c.privkey == nil {
		return nil, errors.New("httpv: No private key set!")
	}

	// Munge req and resp together
	entire_conv, err := munge(c.req, c.resp)
	if err != nil {
		return nil, err
	}

	// Hash
	hash := sha256.Sum256(entire_conv)
	c.hash = hash[:]
	if c.debug {
		fmt.Printf("%s\n", entire_conv)
	}

	// Sign
	sig, err := textSig(*c.privkey, c.hash)
	if err != nil {
		return nil, err
	}

	// Add Signature
	c.resp.Header.Add("Httpv-Sig", sig)

	return c.resp, nil
}

// Takes a conversation and verifies the signature against the public key
// provided to c.
func (c *Conversation) Verify() (bool, error) {
	if c.pubkey == nil {
		return false, errors.New("httpv: No publickey set!")
	}

	// Pull sig out of Response
	sig, err := getSig(c.resp)
	if err != nil {
		return false, err
	}

	// Remove sig from response
	c.resp.Header.Del("Httpv-Sig")

	// Munge conversation
	entire_conv, err := munge(c.req, c.resp)
	if err != nil {
		return false, err
	}
	if c.debug {
		fmt.Printf("%s\n", entire_conv)
	}

	// Hash
	hash := sha256.Sum256(entire_conv)
	c.hash = hash[:]

	// Verifiy the Sig
	v := sig.Verify(c.hash, c.pubkey)
	return v, nil
}

func (c *Conversation) ReadFromFile() {

}

func (c *Conversation) WriteToFile() {

}
