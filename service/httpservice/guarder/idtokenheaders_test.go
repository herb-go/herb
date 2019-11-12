package guarder

import (
	"net/http"
	"testing"

	"github.com/herb-go/herb/service"
)

func TestIDTokenHeaders(t *testing.T) {
	req, err := http.NewRequest("POST", "http://127.0,0,1", nil)
	if err != nil {
		t.Fatal(err)
	}
	c := &service.ConfigMap{}
	c.Set("IDHeader", "id")
	d, err := NewMapperDriver("header", c, "")
	if err != nil {
		t.Fatal(err)
	}
	p, err := d.ReadParamsFromRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if p.ID() != "" || p.Credential() != "" {
		t.Fatal(*p)
	}
	p = NewParams()
	p.SetID("testid")
	p.SetCredential("teestoken")
	d.WriteParamsToRequest(req, p)
	p, err = d.ReadParamsFromRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if p.ID() != "testid" || p.Credential() != "teestoken" {
		t.Fatal(*p)
	}
}
