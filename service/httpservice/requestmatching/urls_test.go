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
	if !mustMatch(p, r) {
		t.Fatal(p)
	}
	p.Add("127.0.0.2/")
	if mustMatch(p, r) {
		t.Fatal(p)
	}
	p.Add("127.0.0.2/path?")
	if mustMatch(p, r) {
		t.Fatal(p)
	}
	p.Add("127.0.0.2/path")
	if mustMatch(p, r) {
		t.Fatal(p)
	}
	p.Add("/path")
	if !mustMatch(p, r) {
		t.Fatal(p)
	}
	p = NewPaths()
	if !mustMatch(p, r) {
		t.Fatal(p)
	}
	p.Add("127.0.0.2/path")
	if mustMatch(p, r) {
		t.Fatal(p)
	}
	p.Add("127.0.0.1/path")
	if !mustMatch(p, r) {
		t.Fatal(p)
	}

}

func TestPrefixs(t *testing.T) {
	var r *http.Request
	var p *Prefixs
	r, _ = http.NewRequest("POST", "http://127.0.0.1/path?a=123", nil)
	p = NewPrefixs()
	if !mustMatch(p, r) {
		t.Fatal(p)
	}
	p.Add("127.0.0.2/")
	if mustMatch(p, r) {
		t.Fatal(p)
	}
	p.Add("127.0.0.2/path?")
	if mustMatch(p, r) {
		t.Fatal(p)
	}
	p.Add("127.0.0.2/path")
	if mustMatch(p, r) {
		t.Fatal(p)
	}
	p.Add("/path")
	if !mustMatch(p, r) {
		t.Fatal(p)
	}
	p = NewPrefixs()
	if !mustMatch(p, r) {
		t.Fatal(p)
	}
	p.Add("127.0.0.2/path")
	if mustMatch(p, r) {
		t.Fatal(p)
	}
	p.Add("127.0.0.1/path")
	if !mustMatch(p, r) {
		t.Fatal(p)
	}
	p = NewPrefixs()
	if !mustMatch(p, r) {
		t.Fatal(p)
	}
	p.Add("127.0.0.2/path")
	if mustMatch(p, r) {
		t.Fatal(p)
	}
	p.Add("127.0.0.1/")
	if !mustMatch(p, r) {
		t.Fatal(p)
	}

}
