package main

import (
	"log"
	"os"

	"github.com/chzar/stratum/v2/internal/server"
)

func main() {
	c, err := server.LoadServerConfigFromFile(os.Args[1])
	log.Default().Printf("Loading %s", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	proxy, _ := server.NewServer(c)
	println("Starting Server...")
	log.Fatal(server.ListenAndServeTLS(":9443", *c.CACert, proxy))
}
