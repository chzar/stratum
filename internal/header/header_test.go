package header

import (
	"bytes"
	"testing"
)

func TestWrite(t *testing.T) {
	b := bytes.NewBuffer(make([]byte, 0))
	h := New(60)
	h.WriteIO(b)

	h2, err := NewFromIO(b)
	if err != nil {
		t.Fatal("Failed to decode header")
	}
	t.Log(h2)
}
