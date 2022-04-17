package header

import (
	"bytes"
	"testing"
)

func TestWriteIO(t *testing.T) {
	b := bytes.NewBuffer(make([]byte, 0))
	h := New(60)
	_, err := h.WriteIO(b)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewFromIO(t *testing.T) {
	b := bytes.NewBuffer(make([]byte, 0))
	h1 := New(60)
	_, err := h1.WriteIO(b)
	if err != nil {
		t.Fatal(err)
	}

	h2, err := NewFromIO(b)
	if err != nil {
		t.Fatal("Failed to decode header")
	}

	if !h1.Expiry.Equal(h2.Expiry) {
		t.Fail()
		t.Log(h1.Expiry)
		t.Log(h2.Expiry)
		t.Log("Expiry times do not match!")
	}
}
