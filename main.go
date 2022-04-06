package main

import (
	"log"
	"os"

	"github.com/chzar/stratum/v2/internal"
)

func main() {
	c, err := internal.LoadServerConfigFromFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	proxy, _ := internal.BuildServer(c)
	println("Starting Server...")
	log.Fatal(internal.ListenAndServeTLS(":9443", *c.CACert, proxy))
}
