package internal

import (
	"net/http"
	"os"

	gp "github.com/elazarl/goproxy"
)

func responseHandler(resp *http.Response, ctx *gp.ProxyCtx) *http.Response {
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
}
