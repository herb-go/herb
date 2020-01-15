package requestmatching

import (
	"net/http"
	"testing"
)

func TestMethods(t *testing.T) {
	var r *http.Request
	r, _ = http.NewRequest("POST", "http://127.0.0.1/", nil)
	m := NewMethods()
	if !mustMatch(m, r) {
		t.Fatal(m)
	}
	m.Add("GET")
	if mustMatch(m, r) {
		t.Fatal(m)
	}
	m.Add("Post")
	if !mustMatch(m, r) {
		t.Fatal(m)
	}
}
