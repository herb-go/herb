package profile

import (
	"strings"
	"testing"
)

func TestProfile(t *testing.T) {
	p := NewProfile()
	if p.Load("name") != "" {
		t.Fatal(p)
	}
	p.
		With("name", "test").
		With("name", "test2").
		With("email", "a@b.c")
	data := p.Data()
	if len(data) != 3 {
		t.Fatal(data)
	}
	result := p.Load("email")
	if result != "a@b.c" {
		t.Fatal(result)
	}
	result = p.Load("name")
	if result != "test2" {
		t.Fatal(result)
	}
	result = p.Load("notexists")
	if result != "" {
		t.Fatal(result)
	}
	p2 := p.Clone()
	if p2 == p {
		t.Fatal()
	}
	result = p.LoadAny()
	if result != "" {
		t.Fatal(result)
	}
	result = p.LoadAny("notexist")
	if result != "" {
		t.Fatal(result)
	}
	result = p.LoadAny("notexist", "name", "email")
	if result != "test2" {
		t.Fatal(result)
	}
	result = p.LoadAny("email", "name", "notexist")
	if result != "a@b.c" {
		t.Fatal(result)
	}
	p.With("newfield", "newvalue")

	filtered := p.Filter("name", "email")
	if len(filtered.Data()) != 3 || filtered.Load("name") != p.Load("name") || filtered.Load("email") != p.Load("email") || filtered.Load("newfield") == p.Load("newfield") {
		t.Fatal(filtered)
	}
	pnew := NewProfile().With("name", "namenew")
	pnew2 := NewProfile().With("name", "namenew2")
	p.Chain(pnew, pnew2)
	if len(p.Data()) != 6 || len(p.Filter("name").Data()) != 4 {
		t.Fatal(p)
	}
	if strings.Join(p.LoadAllByName("name"), ",") != "test,test2,namenew,namenew2" {
		t.Fatal(p)
	}
}
