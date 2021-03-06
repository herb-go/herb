package requestmatching

import (
	"net/http"
	"testing"
)

func TestPaths(t *testing.T) {
	var r *http.Request
	var p *Paths
	r, _ = http.NewRequest("POST", "http://127.0.0.1/path?a=123", nil)
	p = NewPaths()
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
	p = NewPaths()
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
	p = NewPaths()
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

}
