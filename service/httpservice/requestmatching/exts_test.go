package requestmatching

import (
	"net/http"
	"testing"
)

func TestExts(t *testing.T) {
	var r *http.Request
	r, _ = http.NewRequest("POST", "http://127.0.0.1/1.html", nil)
	e := NewExts()
	if !MustMatch(r, e) {
		t.Fatal(e)
	}
	e.Add("")
	if MustMatch(r, e) {
		t.Fatal(e)
	}
	e.Add(".html")
	if !MustMatch(r, e) {
		t.Fatal(e)
	}
	r, _ = http.NewRequest("POST", "http://127.0.0.1/1.HTml", nil)
	e = NewExts()
	e.Add(".htML")
	if !MustMatch(r, e) {
		t.Fatal(e)
	}
}
