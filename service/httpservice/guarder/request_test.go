package guarder_test

import (
	"net/http"
	"testing"

	"github.com/herb-go/herbconfig/loader"
	_ "github.com/herb-go/herbconfig/loader/drivers/jsonconfig"

	"github.com/herb-go/herb/service/httpservice/guarder"
)

func newConfig() *guarder.DriverConfig {
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
	m := guarder.NewDriverConfig()
	err := loader.LoadConfig("json", []byte(config), m)
	if err != nil {
		panic(err)
	}
	return m
}
func TestRequest(t *testing.T) {
	var err error
	c := newConfig()
	g := guarder.NewGuarder()
	err = g.Init(c)
	v := guarder.NewVisitor()
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
	err = v.InitRequest(req)
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
