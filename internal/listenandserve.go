package internal

import (
	"crypto/tls"
	"net"
	"net/http"
)

func ListenAndServeTLS(addr string, cert tls.Certificate,
	handler http.Handler) error {

	if addr == "" {
		addr = "0.0.0.0:443"
	}

	server := &http.Server{Addr: addr, Handler: handler}

	server.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	defer ln.Close()

	return server.ServeTLS(ln, "", "")
}
