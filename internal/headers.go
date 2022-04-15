package internal

import (
	"bytes"
	"encoding/gob"
	"time"
)

type header struct {
	status   int8
	expiry   time.Time
	metadata []byte
}

func NewHeader(TTL uint) *header {
	return &header{
		0,
		time.Now().Add(time.Duration(1000000000 * TTL)),
		make([]byte, 0),
	}
}

func (h *header) GetStatus() int8 {
	return h.status
}

func (h *header) SetStatus(s int8) {
	h.status = s
}

func (h *header) Serialize() []byte {
	b := &bytes.Buffer{}
	gob.NewEncoder(b).Encode(&h)
	return b.Bytes()
}
