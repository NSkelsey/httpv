package main

import (
	"fmt"
	"log"

	"github.com/NSkelsey/httpv"
	"github.com/NSkelsey/net/http"
)

func main() {

	url := "httpv://localhost:8000/helloworld"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Httpv-Ver", httpv.HttpvVer)

	_, pubkey, _ := httpv.FakeKey()
	convo := httpv.NewConversation("localhost:8000", pubkey, nil)
	if err = convo.AddRequest(*req); err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if err = convo.AddResponse(*resp); err != nil {
		log.Fatal(err)
	}

	verified, err := convo.Verify()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Did We verify? %v\n", verified)
}
