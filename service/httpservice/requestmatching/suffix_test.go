package requestmatching

import (
	"net/http"
	"testing"
)

func TestSuffixs(t *testing.T) {
	var r *http.Request
	var p *Suffixs
	r, _ = http.NewRequest("POST", "http://127.0.0.1/url/path?a=123", nil)
	p = NewSuffixs()
	if !MustMatch(r, p) {
		t.Fatal(p)
	}
	p.Add("127.0.0.2/")
	if MustMatch(r, p) {
		t.Fatal(p)
	}
	p.Add("127.0.0.2/path?")
	if MustMatch(r, p) {
		t.Fatal(p)
	}
	p.Add("127.0.0.2/path")
	if MustMatch(r, p) {
		t.Fatal(p)
	}
	p.Add("/path")
	if !MustMatch(r, p) {
		t.Fatal(p)
	}
	p = NewSuffixs()
	if !MustMatch(r, p) {
		t.Fatal(p)
	}
	p.Add("127.0.0.2/path")
	if MustMatch(r, p) {
		t.Fatal(p)
	}
	p.Add("127.0.0.1/path")
	if !MustMatch(r, p) {
		t.Fatal(p)
	}
	p = NewSuffixs()
	if !MustMatch(r, p) {
		t.Fatal(p)
	}
	p.Add("127.0.0.2/path")
	if MustMatch(r, p) {
		t.Fatal(p)
	}
	p.Add("127.0.0.1/Path")
	if !MustMatch(r, p) {
		t.Fatal(p)
	}

	p = NewSuffixs()
	if !MustMatch(r, p) {
		t.Fatal(p)
	}
	p.Add("127.0.0.2/path")
	if MustMatch(r, p) {
		t.Fatal(p)
	}
	p.Add("127.0.0.1/")
	if !MustMatch(r, p) {
		t.Fatal(p)
	}

}
