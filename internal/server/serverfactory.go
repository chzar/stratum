package server

import (
	"net/http"

	gp "github.com/elazarl/goproxy"
)

func NewServer(c *ServerConfig) (*gp.ProxyHttpServer, error) {

	var StratumAlwaysMitm gp.FuncHttpsHandler = func(host string, ctx *gp.ProxyCtx) (*gp.ConnectAction, string) {
		return &gp.ConnectAction{Action: gp.ConnectMitm, TLSConfig: gp.TLSConfigFromCA(c.CACert)}, host
	}

	proxy := gp.NewProxyHttpServer()

	for _, h := range c.Hosts {
		proxy.OnRequest(gp.DstHostIs(h)).HandleConnect(StratumAlwaysMitm)
	}

	for _, re := range c.UrlPatterns {
		proxy.OnRequest(gp.UrlMatches(re)).DoFunc(func(req *http.Request, ctx *gp.ProxyCtx) (*http.Request, *http.Response) {
			ctx.UserData = re.FindStringSubmatch(req.URL.String())[1]
			return req, nil
		})
	}

	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *gp.ProxyCtx) *http.Response {
		return responseHandler(resp, ctx, c)
	})
	proxy.Verbose = true
	return proxy, nil
}
