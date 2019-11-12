package guarder

import (
	"testing"

	"github.com/herb-go/herb/service"
)

func TestDummy(t *testing.T) {
	c := &service.ConfigMap{}
	c.Set("ID", "testid")
	d, err := NewIdentifierDriver("dummy", c, "")
	if err != nil {
		t.Fatal(err)
	}
	p := NewParams()
	id, err := d.IdentifyParams(p)
	if err != nil {
		t.Fatal(err)
	}
	if id != "testid" {
		t.Fatal(id)
	}

}

func TestDummyDefault(t *testing.T) {
	c := &service.ConfigMap{}
	c.Set("ID", "testid")
	d, err := NewIdentifierDriver("", c, "")
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

}
