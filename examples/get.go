package main

import (
	"fmt"
	"log"

	"github.com/NSkelsey/httpv"
)

func main() {

	url := "httpsv://localhost:7999/helloworld"
	_, pubkey, _ := httpv.FakeKey()
	convo, err := httpv.Get(url, pubkey)
	if err != nil {
		log.Fatal(err)
	}

	v, err := convo.Verify()
	fmt.Printf("Result: %v, %v\n=====\n%s\n", v, err, convo.RespBody())
}
