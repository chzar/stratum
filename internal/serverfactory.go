package internal

import (
	"net/http"
	"os"

	gp "github.com/elazarl/goproxy"
)

func BuildServer(c *ServerConfig) (*gp.ProxyHttpServer, error) {

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
		// short circuit non 200 codes
		if resp == nil || resp.StatusCode != 200 || ctx.UserData == nil {
			return resp
		}

		if ctx.UserData.(string) != "" {
			filename := "/tmp/" + ctx.UserData.(string)
			f, err := os.Open(filename)

			// file does not exist, so copy the resp as we write it back to the client
			if err != nil {
				resp.Body = NewTeeReadCloser(resp.Body, NewFileStream(filename))
			} else { // file exists in cache
				ctx.Logf("READING FROM CACHE: %s", filename)
				resp.Body = f
			}
		}
		return resp
	})
	proxy.Verbose = true
	return proxy, nil
}
