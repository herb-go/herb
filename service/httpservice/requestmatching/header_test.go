package requestmatching

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestHeaders(t *testing.T) {
	r, _ := http.NewRequest("POST", "http://127.0.0.1/", nil)
	r.Header.Add("Testfield", "12345")
	h := NewHeaders()
	err := h.Add("12345")
	if !strings.Contains(err.Error(), "12345") || errors.Unwrap(err) != ErrHeaderNotValidated {
		t.Fatal(err)
	}
	h = NewHeaders()
	if !MustMatch(r, h) {
		t.Fatal(h)
	}
	h.Add("Testfield:2345")
	if MustMatch(r, h) {
		t.Fatal(h)
	}
	h.Add("testfield : 2345")
	if MustMatch(r, h) {
		t.Fatal(h)
	}
	h.Add("testfield : 12345")
	if !MustMatch(r, h) {
		t.Fatal(h)
	}
}
