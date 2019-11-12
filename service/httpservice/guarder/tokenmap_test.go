package guarder

import (
	"testing"

	"github.com/herb-go/herb/service"
)

func TestTokenMap(t *testing.T) {
	c := &service.ConfigMap{}
	c.Set("Tokens", map[string]string{"testid": "testtoken"})
	d, err := NewIdentifierDriver("tokenmap", c, "")
	if err != nil {
		t.Fatal(err)
	}
	p := NewParams()
	id, err := d.IdentifyParams(p)
	if err != nil {
		t.Fatal(err)
	}
	if id != "" {
		t.Fatal(id)
	}
	p.SetID("TestID")
	id, err = d.IdentifyParams(p)
	if err != nil {
		t.Fatal(err)
	}
	if id != "" {
		t.Fatal(id)
	}
	p.SetID("testid")
	id, err = d.IdentifyParams(p)
	if err != nil {
		t.Fatal(err)
	}
	if id != "" {
		t.Fatal(id)
	}
	p.SetCredential("testtoken")
	id, err = d.IdentifyParams(p)
	if err != nil {
		t.Fatal(err)
	}
	if id != "testid" {
		t.Fatal(id)
	}
	c.Set("ToLower", true)
	d, err = NewIdentifierDriver("tokenmap", c, "")
	if err != nil {
		t.Fatal(err)
	}
	p.SetID("TestID")
	id, err = d.IdentifyParams(p)
	if err != nil {
		t.Fatal(err)
	}
	if id != "testid" {
		t.Fatal(id)
	}

}
