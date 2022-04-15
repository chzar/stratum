package internal

import (
	"encoding/binary"
	"net/http"
	"os"

	gp "github.com/elazarl/goproxy"
)

const TTL = 6000

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
			// create a header
			h := NewHeader(TTL)
			// create a filestream
			fs := NewFileStream(filename)
			// serialize the header
			hBytes := h.Serialize()
			// write length of header as a little endian int64 to beginning of file
			hLenBytes := make([]byte, 8)
			binary.LittleEndian.PutUint64(hLenBytes, uint64(len(hBytes)))
			fs.Write(hLenBytes)
			// write the header to the response
			fs.Write(hBytes)
			resp.Body = NewTeeReadCloser(resp.Body, NewFileStream(filename))
		} else { // file exists in cache
			ctx.Logf("READING FROM CACHE: %s", filename)
			// read the first 8 bytes to get length of header
			hLenBytes := make([]byte, 8)
			f.Read(hLenBytes)
			hLen := binary.LittleEndian.Uint64(hLenBytes)
			hBytes := make([]byte, hLen)
			f.Read(hBytes)
			resp.Body = f
		}
	}
	return resp
}
