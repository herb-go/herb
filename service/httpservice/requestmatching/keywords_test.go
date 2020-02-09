package requestmatching

import (
	"net/http"
	"testing"
)

func TestKeywords(t *testing.T) {
	var r *http.Request
	r, _ = http.NewRequest("POST", "http://127.0.0.1/1.html", nil)
	k := NewKeywords()
	if !MustMatch(r, k) {
		t.Fatal(k)
	}
	k.Add("notexists")
	if MustMatch(r, k) {
		t.Fatal(k)
	}
	k.Add(".html")
	if !MustMatch(r, k) {
		t.Fatal(k)
	}
	r, _ = http.NewRequest("POST", "http://127.0.0.1/1.HTml", nil)
	k = NewKeywords()
	k.Add(".htML")
	if !MustMatch(r, k) {
		t.Fatal(k)
	}
}
