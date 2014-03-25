package main

import (
	"./checkversion"
	"flag"
	"log"
	"net/http"
	"strconv"
)

var port *int = flag.Int("port", 23456, "port to listen.")

func main() {
	flag.Parse()

	server := checkversion.NewServer("/checkversion")
	go server.Listen()
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
