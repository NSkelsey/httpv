// This file was cannibalized from the Go Authors
// Specifically: http://golang.org/src/pkg/net/http/httputil/reverseproxy.go
// All credit and copyright claims go to them.

package httpv

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/conformal/btcec"
)

// An HTTPv reverse proxy handler
type RevProxy struct {
	// The transport used to perform proxy requests.
	transport http.RoundTripper
	target    *url.URL
	privkey   *btcec.PrivateKey
	host      string
}

func NewReverseProxy(host, target string, privkey *btcec.PrivateKey) (*RevProxy, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	rp := &RevProxy{
		transport: http.DefaultTransport,
		target:    u,
		privkey:   privkey,
		host:      host,
	}
	return rp, nil
}

// Really Serves an HTTPv connection that is being proxied to the RevProxy's
// target. All paths are copied from straight from the incoming request to
// the proxied one.
func (rp *RevProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	u := rp.target
	u.Path = req.URL.Path
	if req.Method != "GET" {
		http.Error(rw, "Bad!", http.StatusMethodNotAllowed)
		return
	}

	privkey, _, _ := FakeKey()

	// Handle the httpv request logic
	convo := NewConversation(rp.host, nil, privkey)

	if err := convo.AddRequest(*req); err != nil {
		log.Printf("%s\n", err)
		http.Error(rw, err.Error(), 700)
		return
	}

	outreq, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Printf("%s\n", err)
		http.Error(rw, err.Error(), 500)
		return
	}
	outreq.Close = false

	copyHeaders(outreq.Header, req.Header)

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		// If we aren't the first proxy retain prior
		// X-Forwarded-For information as a comma+space
		// separated list and fold multiple headers into one.
		if prior, ok := outreq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outreq.Header.Set("X-Forwarded-For", clientIP)
	}

	proxyResp, err := http.DefaultClient.Do(outreq)
	if err != nil {
		log.Printf("httpv: proxy error: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	proxyResp.Close = false
	defer proxyResp.Body.Close()

	if proxyResp.StatusCode != 200 {
		http.Error(rw, proxyResp.Status, proxyResp.StatusCode)
		return
	}

	// Now that we have the proxy's response we can finish the httpv response.
	if err = convo.AddResponse(*proxyResp); err != nil {
		log.Printf("%s\n", err)
		http.Error(rw, err.Error(), 700)
		return
	}

	httpvResp, err := convo.EmitResponse()
	if err != nil {
		log.Printf("%s\n", err)
		http.Error(rw, err.Error(), 700)
		return
	}
	httpvResp.Close = false

	copyHeaders(rw.Header(), httpvResp.Header)

	if err != nil {
		log.Fatal(err)
	}
	rw.Write(convo.RespBody())
}

func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
