package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NSkelsey/httpv"
)

func main() {
	proxyfor := "http://localhost:8000"
	iface := "localhost:7999"

	privkey, _, _ := httpv.FakeKey()
	rp, err := httpv.NewReverseProxy(iface, proxyfor, privkey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Starting server...\n"+
		"Proxying to: %s\n"+
		"Listening at: %s\n",
		proxyfor, iface)
	log.Fatal(http.ListenAndServeTLS(iface, "cert.pem", "key.pem", rp))
}
