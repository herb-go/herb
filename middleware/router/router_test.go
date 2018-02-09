package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herb-go/herb/middleware"
)

func TestParams(t *testing.T) {
	app := middleware.New()
	middlewareSetParam := func(name, value string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			params := GetParams(r)
			params.Set(name, value)
			next(w, r)
		}
	}
	action := func(w http.ResponseWriter, r *http.Request) {
		params := GetParams(r)
		w.Header().Set("test1", params.Get("test1"))
		w.Header().Set("test2", params.Get("test2"))
		w.Header().Set("testNotExists", params.Get("testNotExists"))
		bs, err := json.Marshal(params)
		if err != nil {
			panic(err)
		}
		_, err = w.Write(bs)
		if err != nil {
			panic(err)
		}
	}
	app.Use(
		middlewareSetParam("test1", "testa"),
		middlewareSetParam("test2", "test2"),
		middlewareSetParam("test1", "test1"),
		middlewareSetParam("test3", "test3"),
	).HandleFunc(action)
	server := httptest.NewServer(app)
	defer server.Close()
	resp, err := http.Get(server.URL)
	if err != nil {
		panic(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	result := Params{}
	err = json.Unmarshal(content, &result)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 4 {
		t.Error(len(result))
	}
	resp.Body.Close()
	if err != nil {
		panic(err)
	}
	if resp.Header.Get("test1") != "test1" {
		t.Error(resp.Header.Get("test1"))
	}
	if resp.Header.Get("test2") != "test2" {
		t.Error(resp.Header.Get("test2"))
	}
	if resp.Header.Get("testNotExists") != "" {
		t.Error(resp.Header.Get("testNotExists"))
	}

}

func TestNilParams(t *testing.T) {
	var p Params
	if p.Get("test") != "" {
		t.Error(p.Get("test"))
	}
	p.Set("test", "test")
	if p.Get("test") != "test" {
		t.Error(p.Get("test"))
	}
	p = Params{}
	if p.Get("test") != "" {
		t.Error(p.Get("test"))
	}
	p.Set("test", "test")
	if p.Get("test") != "test" {
		t.Error(p.Get("test"))
	}
}
