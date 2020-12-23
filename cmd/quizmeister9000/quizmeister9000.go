package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/thijzert/quizmeister9000/qm9k"
)

func main() {
	var config qm9k.Config
	var addr string

	flag.StringVar(&addr, "addr", "localhost:20598", "http service address")
	flag.BoolVar(&config.SecureCookies, "secure-cookies", false, "Set the 'secure' flag on all cookies")
	flag.Parse()

	server, err := qm9k.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting web server on '%s'", addr)
	err = http.ListenAndServe(addr, server)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
