package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/pires/go-proxyproto"
)

var (
	targetEndpoint = os.Getenv("TARGET_ENDPOINT")
	addr           = os.Getenv("ADDR")
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
		log.Println("request received")

		res, err := http.Get(targetEndpoint)
		if err != nil {
			rw.WriteHeader(http.StatusBadGateway)
			fmt.Printf("failed to hit upstream service: %s", err)
			return
		}

		resBytes, _ := ioutil.ReadAll(res.Body)
		fmt.Printf("upstream service responded with %s", string(resBytes))

		rw.WriteHeader(http.StatusOK)
		fmt.Fprintf(rw, `{"remoteAddr":"%s","upstreamResponse":%s}`, r.RemoteAddr, string(resBytes))
	}))

	server := http.Server{Addr: addr, Handler: mux}
	log.Fatalln(server.Serve(ppls))
}

func main() {
	httpServer()
}
