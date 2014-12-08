package main

import (
	"bufio"

	"io"
	"log"
	"net"

	"github.com/NSkelsey/httpv"
	"github.com/NSkelsey/net/http"
)

func writeErr(w *bufio.Writer, err error) {

	log.Println(err)
	errResp := http.Response{
		Status:     "500 Internal Error",
		StatusCode: 500,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     http.Header{},
	}
	errResp.Write(w)

	w.Flush()
}

func handleConn(c net.Conn) {
	defer c.Close()
	log.Println("Conn recieved!")
	lr := io.LimitReader(c, (1<<63)-1)
	reader := bufio.NewReader(lr)
	writer := bufio.NewWriterSize(c, 4<<10)
	w := writer

	var req *http.Request
	var err error
	req, err = http.ReadRequest(reader)
	if err != nil {
		log.Fatal(err)
	}

	resp := http.Response{
		Status:        "204 No Content",
		StatusCode:    204,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        http.Header{},
		Request:       req,
		ContentLength: 0,
	}

	privkey, _, _ := httpv.FakeKey()
	convo := httpv.NewConversation("localhost:8000", nil, privkey)

	if err = convo.AddRequest(*req); err != nil {
		writeErr(w, err)
		return
	}

	if err = convo.AddResponse(resp); err != nil {
		writeErr(w, err)
		return
	}

	signedresp, err := convo.EmitResponse()
	if err != nil {
		writeErr(w, err)
		return
	}

	err = signedresp.Write(writer)
	if err != nil {
		writeErr(w, err)
		return
	}
	writer.Flush()

}

func main() {
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConn(c)
	}
}
