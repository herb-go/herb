package guarder

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestTokenMap(t *testing.T) {
	var err error

	buf := bytes.NewBuffer(nil)
	err = json.NewEncoder(buf).Encode(TokenMap{
		Tokens: map[string]string{"testid": "testtoken"},
	})
	if err != nil {
		t.Fatal(err)
	}

	d, err := NewIdentifierDriver("tokenmap", json.NewDecoder(buf).Decode)
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

	err = json.NewEncoder(buf).Encode(TokenMap{
		Tokens:  map[string]string{"testid": "testtoken"},
		ToLower: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	d, err = NewIdentifierDriver("tokenmap", json.NewDecoder(buf).Decode)
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
