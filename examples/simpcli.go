package main

import (
	"fmt"
	"log"

	"github.com/NSkelsey/httpv"
)

func main() {

	url := "httpv://localhost:8000/helloworld"

	_, pubkey, _ := httpv.FakeKey()
	convo, err := httpv.Get(url, pubkey)
	if err != nil {
		log.Fatal(err)
	}

	verified, err := convo.Verify()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Did We verify? %v\n", verified)
}
