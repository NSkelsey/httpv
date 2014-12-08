package httpv

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"strconv"
	"time"

	"github.com/NSkelsey/net/http"
	"github.com/conformal/btcec"
)

var (
	ErrReqNoVerHead = errors.New("httpv: Client does not have a HTTPV-Ver header")
	ErrReqMethod    = errors.New("httpv: Request must be a GET to sign")
)

func respHeaders(header http.Header) (http.Header, error) {

	newmap := make(map[string][]string)
	// copy old headers
	for k, v := range map[string][]string(header) {
		newmap[k] = v
	}
	newhead := http.Header(newmap)

	rb := make([]byte, 32)
	_, err := rand.Read(rb)
	if err != nil {
		return newhead, err
	}

	newhead.Add("Httpv-Ver", HttpvVer)
	ts := strconv.Itoa(int(time.Now().Unix()))
	newhead.Add("Httpv-Ts", ts)

	rs := base64.StdEncoding.EncodeToString(rb)
	newhead.Add("Httpv-N", rs)

	return newhead, nil
}

// Copies a set of headers over from req.header into a new header map
func reqHeaders() http.Header {

	newmap := make(map[string][]string)
	newhead := http.Header(newmap)

	newhead.Add("Httpv-Ver", HttpvVer)

	return newhead
}

// Throws an error based if the Request does not conform to an httpv request.
func enforceReq(req http.Request) error {
	if req.Method != "GET" {
		return ErrReqMethod
	}
	if v := req.Header.Get("Httpv-Ver"); v != HttpvVer {
		return ErrReqNoVerHead
	}

	return nil
}

// munges the request and the response together
func munge(req *http.Request, resp *http.Response) ([]byte, error) {
	empt := []byte{}

	buf := bytes.NewBuffer([]byte{})

	if err := req.Write(buf); err != nil {
		return empt, err
	}

	if err := resp.Write(buf); err != nil {
		return empt, err
	}

	return buf.Bytes(), nil
}

func textSig(privkey btcec.PrivateKey, hash []byte) (string, error) {
	sig, err := privkey.Sign(hash[:])
	if err != nil {
		return "", err
	}

	// encode as base64 and attach to resp
	sigB := sig.Serialize()[:]
	enc_sig := base64.StdEncoding.EncodeToString(sigB)

	return enc_sig, nil
}
