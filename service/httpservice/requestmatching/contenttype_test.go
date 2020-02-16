package requestmatching

import (
	"net/http"
	"testing"
)

func TestContenttypes(t *testing.T) {
	var r *http.Request
	r, _ = http.NewRequest("POST", "http://127.0.0.1/1.html", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ct := NewContentTypes()
	if !MustMatch(r, ct) {
		t.Fatal(ct)
	}
	ct.Add("")
	if MustMatch(r, ct) {
		t.Fatal(ct)
	}
	ct.Add("application/x-www-form-urlencoded")
	if !MustMatch(r, ct) {
		t.Fatal(ct)
	}
	r, _ = http.NewRequest("POST", "http://127.0.0.1/1.HTml", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	ct = NewContentTypes()
	ct.Add("application/x-www-form-urlencoded")
	if !MustMatch(r, ct) {
		t.Fatal(ct)
	}
}
