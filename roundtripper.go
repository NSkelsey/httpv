package httpv

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/conformal/btcec"
)

var (
	ErrBadVerify error = errors.New("httpv: Verification of the conversation failed!")
)

type httpvTransport struct {
}

type httpsvTransport struct {
}

func (t *httpvTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (t *httpsvTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "https"
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewTransport() http.RoundTripper {
	trans := &http.Transport{}
	trans.RegisterProtocol("httpv", &httpvTransport{})
	trans.RegisterProtocol("httpsv", &httpsvTransport{})
	return trans
}

func Get(urlStr string, pubkey *btcec.PublicKey) (*Conversation, error) {

	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "httpv" && u.Scheme != "httpvs" {
		return nil, errors.New("httpv: Cannot use httpv.Get with non httpv scheme.")
	}

	client := http.Client{
		Transport: NewTransport(),
	}

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Httpv-Ver", "0.1")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	convo := NewConversation(u.Host, pubkey, nil)
	if err = convo.AddRequest(*req); err != nil {
		return nil, err
	}

	if err = convo.AddResponseAssert(*resp, true); err != nil {
		return nil, err
	}
	resp.Body.Close()

	return &convo, nil
}
