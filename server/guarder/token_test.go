package guarder

import (
	"testing"

	"github.com/herb-go/herb/server"
)

func TestToken(t *testing.T) {
	c := &server.ConfigMap{}
	c.Set("Token", "testtoken")
	idDriver, err := NewIdentifierDriver("token", c, "")
	if err != nil {
		t.Fatal(err)
	}
	cDriver, err := NewCredentialDriver("token", c, "")
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
	c.Set("ID", "testid")
	idDriver, err = NewIdentifierDriver("token", c, "")
	if err != nil {
		t.Fatal(err)
	}
	cDriver, err = NewCredentialDriver("token", c, "")
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
