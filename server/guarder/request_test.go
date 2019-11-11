package guarder

import (
	"encoding/json"
	"net/http"
	"testing"
)

func newConfig() *DirverConfigMap {
	config := `{
		"Driver":"token",
		"MapperDriver":"header",
		"Config":{
			"IDHeader":"id",
			"TokenHeader":"token",
			"ID":"testid",
			"Token":"testtoken"
		}
		}`
	m := NewDriverConfigMap()
	err := json.Unmarshal([]byte(config), m)
	if err != nil {
		panic(err)
	}
	return m
}
func TestRequest(t *testing.T) {
	var err error
	c := newConfig()
	g := NewGuarder()
	err = g.Init(c)
	v := NewVisitor()
	err = v.Init(c)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "http://127.0.0.1/", nil)
	if err != nil {
		t.Fatal(err)
	}
	id, err := g.IdentifyRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if id != "" {
		t.Fatal(id)
	}
	err = v.CredentialRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	id, err = g.IdentifyRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if id != "testid" {
		t.Fatal(id)
	}
}
