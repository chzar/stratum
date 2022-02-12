package internal

import gp "github.com/elazarl/goproxy"

func BuildServer(c *ServerConfig) *gp.ProxyHttpServer {
	return new(gp.ProxyHttpServer)
}
