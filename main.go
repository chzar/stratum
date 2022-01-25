package main

import (
	"log"
	"net/http"
	"os"
	"regexp"

	"stratum/internal"

	"github.com/elazarl/goproxy"
)

func main() {
	re := regexp.MustCompile("production.cloudflare.docker.com:443/registry-v2/docker/registry/v2/blobs/sha256/.+/(.+)/.*")

	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest(goproxy.DstHostIs("production.cloudflare.docker.com:443")).HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest(goproxy.UrlMatches(re)).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		ctx.UserData = re.FindStringSubmatch(req.URL.String())[1]
		return req, nil
	})

	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
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
