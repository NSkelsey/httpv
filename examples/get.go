package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/NSkelsey/httpv"
)

var (
	urlStr   *string = flag.String("url", "", "The URL of the request.")
	verbose  *bool   = flag.Bool("v", false, "The verboisty of the output.")
	toFile   *string = flag.String("tofile", "", "The file to record the bundle in.")
	fromFile *string = flag.String("fromfile", "", "The load a bundle from.")
)

func main() {

	flag.Parse()
	args := flag.Args()

	var u *url.URL
	var err error
	if len(args) > 0 {
		u, err = url.Parse(args[0])
	} else {
		u, err = url.Parse(*urlStr)
	}
	if err != nil {
		log.Fatal(err)
	}
	if u.String() == "" {
		fmt.Println("url cannot be empty")
		return
	}

	// TODO handle from and to file.
	_, pubkey, _ := httpv.FakeKey()
	fmt.Printf("Requesting: [%s]\n", u.String())
	convo, err := httpv.Get(u.String(), pubkey)
	if err != nil {
		fmt.Println(err)
		return
	}

	v, err := convo.Verify()
	if v {
		fmt.Printf("The signature verified at: %s\n", time.Now())

		if *verbose {
			b, err := convo.Bundle()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("The Bundle: ===========================================\n%s\n", b)
		}
	} else {
		fmt.Printf("The signature failed verification with: %s\n", err)
	}
}
