package server

import (
	"net/http"
	"os"
	"path"

	"github.com/chzar/stratum/v2/internal/header"
	gp "github.com/elazarl/goproxy"
)

const TTL = 6000

func responseHandler(resp *http.Response, ctx *gp.ProxyCtx, c *ServerConfig) *http.Response {
	// short circuit non 200 codes
	if resp == nil || resp.StatusCode != 200 || ctx.UserData == nil {
		return resp
	}

	if ctx.UserData.(string) != "" {
		filename := path.Join("/tmp/", ctx.UserData.(string))
		f, err := os.Open(filename)

		// file does not exist, so copy the resp as we write it back to the client
		if err != nil {
			// create a header
			h := header.New(TTL)
			// create a filestream
			fs := NewFileStream(filename)
			// write the header
			h.WriteIO(fs)

			resp.Body = NewTeeReadCloser(resp.Body, NewFileStream(filename))
		} else { // file exists in cache
			ctx.Logf("READING FROM CACHE: %s", filename)
			header.NewFromIO(f)
			resp.Body = f
		}
	}
	return resp
}
