package muxrouter

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herb-go/herb/middleware/router"
)

var _ router.Router = New()
var testAction = func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.URL.Path))
}

func catch(f func()) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = r.(error)
		}
	}()
	f()
	return nil
}

func TestMuxRouter(t *testing.T) {
	router := New()
	server := httptest.NewServer(router)
	defer server.Close()
	resp, err := http.DefaultClient.Get(server.URL + "/notexist")
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Fatal(resp)
	}

	router.SetNotFoundHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(405), 405)
	}))

	resp, err = http.DefaultClient.Get(server.URL + "/notexist")
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 405 {
		t.Fatal(resp)
	}
	router.StripPrefix("/prefixed").HandleFunc(testAction)

	err = catch(func() {
		router.StripPrefix("/prefixed/").HandleFunc(testAction)
	})
	if err == nil {
		t.Fatal()
	}

	err = catch(func() {
		router.StripPrefix("").HandleFunc(testAction)
	})
	if err == nil {
		t.Fatal()
	}

	router.Handle("/notprefixed/").HandleFunc(testAction)
	router.HandleHomepage().HandleFunc(testAction)

	resp, err = http.DefaultClient.Get(server.URL + "/prefixed/prefixremoved")
	if err != nil {
		panic(err)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if string(data) != "/prefixremoved" {
		t.Fatal(string(data))
	}
	resp, err = http.DefaultClient.Get(server.URL + "/notprefixed/test")
	if err != nil {
		panic(err)
	}
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if string(data) != "/notprefixed/test" {
		t.Fatal(string(data))
	}
	resp, err = http.DefaultClient.Get(server.URL + "/")
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if string(data) != "/notprefixed/test" {
		t.Fatal(string(data))
	}
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
	err = catch(func() {
		router.HandleHomepage().HandleFunc(testAction)
	})
	if err == nil {
		t.Fatal()
	}

}
