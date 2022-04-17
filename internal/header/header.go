package header

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"time"
)

// Gob requires public fields, consider using an anonymous struct for enc/dec
type header struct {
	Status   int8
	Expiry   time.Time
	Metadata string
}

func New(TTL uint) *header {
	return &header{
		0,
		time.Now().Add(time.Duration(1000000000 * TTL)),
		"foobar",
	}
}

func (h *header) Serialize() []byte {
	b := &bytes.Buffer{}
	gob.NewEncoder(b).Encode(h)
	return b.Bytes()
}

func Deserialize(r io.Reader) (*header, error) {
	h := new(header)

	err := gob.NewDecoder(r).Decode(h)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (h *header) WriteIO(w io.Writer) (int, error) {
	hBytes := h.Serialize()

	// write length of header as a little endian uint64 to beginning of file
	var tbw int

	hLenBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(hLenBytes, uint64(len(hBytes)))

	bw, err := w.Write(hLenBytes)
	if err != nil {
		return bw, err
	}

	tbw = tbw + bw

	// write the header
	bw, err = w.Write(hBytes)
	tbw = tbw + bw
	if err != nil {
		return tbw, err
	}

	return tbw, nil
}

func NewFromIO(r io.Reader) (*header, error) {
	hLenBytes := make([]byte, 8)

	_, err := r.Read(hLenBytes)
	if err != nil {
		return nil, err
	}

	hLen := binary.LittleEndian.Uint64(hLenBytes)
	hBytes := make([]byte, hLen+8)

	// write the header
	_, err = r.Read(hBytes)
	if err != nil {
		return nil, err
	}

	return Deserialize(bytes.NewReader(hBytes))
}

func (h *header) String() string {
	return fmt.Sprintf("Status: %d, Expiry: %s", h.Status, h.Expiry)
}
