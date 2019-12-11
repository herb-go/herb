package guarder

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestDummy(t *testing.T) {
	var err error
	buf := bytes.NewBuffer(nil)
	err = json.NewEncoder(buf).Encode(Dummy{
		ID: "testid",
	})
	if err != nil {
		t.Fatal(err)
	}

	d, err := NewIdentifierDriver("dummy", json.NewDecoder(buf).Decode)
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
	var err error
	buf := bytes.NewBuffer(nil)
	err = json.NewEncoder(buf).Encode(Dummy{
		ID: "testid",
	})
	if err != nil {
		t.Fatal(err)
	}

	d, err := NewIdentifierDriver("", json.NewDecoder(buf).Decode)
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
