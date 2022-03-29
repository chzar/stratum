package main

import (
	"log"

	"github.com/chzar/stratum/v2/internal"
)

func main() {
	c := internal.LoadServerConfig()
	proxy, _ := internal.BuildServer(c)
	println("Starting Server...")
	log.Fatal(internal.ListenAndServeTLS(":9443", *c.CACert, proxy))
}
