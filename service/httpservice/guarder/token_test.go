package guarder

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestToken(t *testing.T) {
	var err error
	buf := bytes.NewBuffer(nil)
	err = json.NewEncoder(buf).Encode(Token{
		Token: "testtoken",
	})
	if err != nil {
		t.Fatal(err)
	}
	idDriver, err := NewIdentifierDriver("token", json.NewDecoder(buf).Decode)
	if err != nil {
		t.Fatal(err)
	}
	err = json.NewEncoder(buf).Encode(Token{
		Token: "testtoken",
	})
	if err != nil {
		t.Fatal(err)
	}

	cDriver, err := NewCredentialDriver("token", json.NewDecoder(buf).Decode)
	if err != nil {
		t.Fatal(err)
	}
	p := NewParams()
	id, err := idDriver.IdentifyParams(p)
	if err != nil {
		t.Fatal(err)
	}
	if id != "" {
		t.Fatal(id)
	}
	p, err = cDriver.CredentialParams()
	if err != nil {
		t.Fatal(err)
	}
	id, err = idDriver.IdentifyParams(p)
	if err != nil {
		t.Fatal(err)
	}
	if id != DefaultStaticID {
		t.Fatal(id)
	}
	err = json.NewEncoder(buf).Encode(Token{
		ID:    "testid",
		Token: "testtoken",
	})
	if err != nil {
		t.Fatal(err)
	}

	idDriver, err = NewIdentifierDriver("token", json.NewDecoder(buf).Decode)
	if err != nil {
		t.Fatal(err)
	}
	err = json.NewEncoder(buf).Encode(Token{
		ID:    "testid",
		Token: "testtoken",
	})
	if err != nil {
		t.Fatal(err)
	}

	cDriver, err = NewCredentialDriver("token", json.NewDecoder(buf).Decode)
	if err != nil {
		t.Fatal(err)
	}

	p, err = cDriver.CredentialParams()
	if err != nil {
		t.Fatal(err)
	}
	id, err = idDriver.IdentifyParams(p)
	if err != nil {
		t.Fatal(err)
	}
	if id != "testid" {
		t.Fatal(id)
	}

}
