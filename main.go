package main

import (
	"log"
	"net/http"

	"github.com/chzar/stratum/v2/internal"
)

func main() {
	c := internal.LoadServerConfig()
	proxy, _ := internal.BuildServer(c)
	log.Fatal(http.ListenAndServeTLS(":9443", "server.crt", "server.key", proxy))
}
