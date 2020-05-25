package middlewarefactory_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herb-go/herb/middleware"
	"github.com/herb-go/herb/middleware/middlewarefactory"
	"github.com/herb-go/herbconfig/loader"
	_ "github.com/herb-go/herbconfig/loader/drivers/jsonconfig"
)

func mustNewLoader(v interface{}) func(v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	data := string(bs)
	fmt.Println(data)
	return loader.NewLoader("json", bs)
}

var emptyConfig = middlewarefactory.ResponseMiddleware{}

var emptyBodyConfig = middlewarefactory.ResponseMiddleware{
	StatusCode: 404,
}
var notfoundmsg = "NotFound"
var fullConfig = middlewarefactory.ResponseMiddleware{
	StatusCode: 404,
	Body:       &notfoundmsg,
	Header: http.Header{
		"notfoundheader": []string{"notfound"},
	},
}

func TestResponse(t *testing.T) {
	var m middleware.Middleware
	var err error
	app := middleware.New()
	app.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		if m == nil {
			w.Header().Add("nil", "true")
			w.Write([]byte("ok"))
			return
		}
		m(w, r, func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
	})
	s := httptest.NewServer(app)
	defer s.Close()
	f := middlewarefactory.NewResponseFactory()
	l := mustNewLoader(emptyConfig)
	resp, err := http.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 || resp.Header.Get("nil") != "true" {
		t.Fatal(resp)
	}
	resp.Body.Close()
	m, err = f(l)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 || resp.Header.Get("nil") == "true" {
		t.Fatal(resp)
	}
	resp.Body.Close()
	m, err = f(mustNewLoader(emptyBodyConfig))
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 404 || string(buf) != http.StatusText(404) {
		t.Fatal(resp)
	}
	resp.Body.Close()
	m, err = f(mustNewLoader(fullConfig))
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 404 || string(buf) != notfoundmsg || resp.Header.Get("notfoundheader") != "notfound" {
		t.Fatal(resp)
	}
	resp.Body.Close()

}
