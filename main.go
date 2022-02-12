package main

import (
	"log"
	"net/http"
	"os"
	"regexp"

	"stratum/internal"

	gp "github.com/elazarl/goproxy"
)

var StratumMitmConnect = gp.ConnectAction{Action: gp.ConnectMitm, TLSConfig: gp.TLSConfigFromCA(&gp.GoproxyCa)}

func main() {
	re := regexp.MustCompile("production.cloudflare.docker.com:443/registry-v2/docker/registry/v2/blobs/sha256/.+/(.+)/.*")

	proxy := gp.NewProxyHttpServer()
	proxy.OnRequest(gp.DstHostIs("production.cloudflare.docker.com:443")).HandleConnect(gp.AlwaysMitm)
	proxy.OnRequest(gp.UrlMatches(re)).DoFunc(func(req *http.Request, ctx *gp.ProxyCtx) (*http.Request, *http.Response) {
		ctx.UserData = re.FindStringSubmatch(req.URL.String())[1]
		return req, nil
	})

	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *gp.ProxyCtx) *http.Response {
		// short circuit non 200 codes
		if resp == nil || resp.StatusCode != 200 || ctx.UserData == nil {
			return resp
		}

		if ctx.UserData.(string) != "" {
			filename := "workdir/" + ctx.UserData.(string)
			f, err := os.Open(filename)
			if err != nil {
				resp.Body = internal.NewTeeReadCloser(resp.Body, internal.NewFileStream(filename))
			} else {
				ctx.Logf("READING FROM CACHE: %s", filename)
				resp.Body = f
			}
		}
		return resp
	})
	proxy.Verbose = true
	log.Fatal(http.ListenAndServeTLS(":9443", "server.crt", "server.key", proxy))
}
