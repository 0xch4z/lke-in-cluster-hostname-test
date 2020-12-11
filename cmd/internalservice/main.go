package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/pires/go-proxyproto"
)

var (
	addr = os.Getenv("ADDR")
)

func init() {
	if addr == "" {
		addr = ":8080"
	}
}

func httpServer() {
	ls, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("could not listen to %q: %q", addr, err.Error())
	}

	ppls := &proxyproto.Listener{
		Listener: ls,
		ValidateHeader: func(header *proxyproto.Header) error {
			fmt.Printf("proxyprotocol info: %#v\n", header)
			return nil
		},
		Policy: func(upstream net.Addr) (proxyproto.Policy, error) {
			fmt.Printf("proxyprotocol upstream: %#v\n", upstream)
			return proxyproto.REQUIRE, nil
		},
	}
	defer ppls.Close()

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		fmt.Fprintf(rw, `{"remoteAddr":"%s", "message":"hello, world"}`, r.RemoteAddr)
	}))

	server := http.Server{Addr: addr, Handler: mux}
	log.Fatalln(server.Serve(ppls))
}

func main() {
	httpServer()
}
