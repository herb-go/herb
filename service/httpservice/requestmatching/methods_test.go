package requestmatching

import (
	"net/http"
	"testing"
)

func TestMethods(t *testing.T) {
	var r *http.Request
	r, _ = http.NewRequest("POST", "http://127.0.0.1/", nil)
	m := NewMethods()
	if !MustMatch(r, m) {
		t.Fatal(m)
	}
	m.Add("GET")
	if MustMatch(r, m) {
		t.Fatal(m)
	}
	m.Add("Post")
	if !MustMatch(r, m) {
		t.Fatal(m)
	}
}
