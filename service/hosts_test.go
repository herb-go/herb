package service

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/herb-go/herb/middleware"
)

func TestHosts(t *testing.T) {
	var successAction = func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		if err != nil {
			panic(err)
		}
	}
	h := Hosts{}
	app := middleware.New()
	app.Use(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		h.ServeMiddleware(w, r, next)
	}).HandleFunc(successAction)

	s := httptest.NewServer(app)
	defer s.Close()
	if !strings.HasPrefix(s.URL, "http://127.0.0.1") {
		t.Fatal(s.URL)
	}
	resp, err := http.DefaultClient.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
	resp.Body.Close()
	h = Hosts{Patterns: []HostPattern{"notexists"}}
	resp, err = http.DefaultClient.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 404 {
		t.Fatal(resp)
	}
	resp.Body.Close()
	resp.Body.Close()
	h = Hosts{Patterns: []HostPattern{""}}
	resp, err = http.DefaultClient.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 404 {
		t.Fatal(resp)
	}
	resp.Body.Close()
	h = Hosts{Patterns: []HostPattern{"", "127.0.0.1"}}
	resp, err = http.DefaultClient.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
	resp.Body.Close()
	h = Hosts{Patterns: []HostPattern{".1.0.1"}}
	resp, err = http.DefaultClient.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 404 {
		t.Fatal(resp)
	}
	resp.Body.Close()
	h = Hosts{Patterns: []HostPattern{".0.1"}}
	resp, err = http.DefaultClient.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 404 {
		t.Fatal(resp)
	}
	resp.Body.Close()
	h = Hosts{Patterns: []HostPattern{".0.0.1"}}
	resp, err = http.DefaultClient.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
	resp.Body.Close()
	h = Hosts{Patterns: []HostPattern{"*.1.1"}}
	resp, err = http.DefaultClient.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 404 {
		t.Fatal(resp)
	}
	resp.Body.Close()
	h = Hosts{Patterns: []HostPattern{"*.0.1"}}
	resp, err = http.DefaultClient.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
	resp.Body.Close()
}
