package misc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herb-go/herb/middleware"
)

func TestIf(t *testing.T) {
	var successAction = func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		if err != nil {
			panic(err)
		}
	}
	condition := true
	appiftrue := middleware.New()
	appiftrue.Use(ErrorIf(true, 400)).HandleFunc(successAction)
	appiffalse := middleware.New()
	appiffalse.Use(ErrorIf(false, 400)).HandleFunc(successAction)
	appwhen := middleware.New(ErrorWhen(func() (bool, error) { return condition, nil }, 400))
	mux := http.NewServeMux()
	mux.Handle("/iftrue", appiftrue)
	mux.Handle("/iffalse", appiffalse)
	mux.Handle("/when", appwhen)
	server := httptest.NewServer(mux)
	defer server.Close()
	resp, err := http.Get(server.URL + "/iftrue")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Error(resp.StatusCode)
	}
	resp, err = http.Get(server.URL + "/iffalse")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode == 400 {
		t.Error(resp.StatusCode)
	}
	resp, err = http.Get(server.URL + "/when")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Error(resp.StatusCode)
	}
	condition = false
	resp, err = http.Get(server.URL + "/when")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode == 400 {
		t.Error(resp.StatusCode)
	}
}

func TestWhen(t *testing.T) {
	var successAction = func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		if err != nil {
			panic(err)
		}
	}
	var m = func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.Header().Set("Condition", "true")
		next(w, r)
	}
	condition := true
	appiftrue := middleware.New()
	appiftrue.Use(MiddlewareIf(true, m)).HandleFunc(successAction)
	appiffalse := middleware.New()
	appiffalse.Use(MiddlewareIf(false, m)).HandleFunc(successAction)
	appwhen := middleware.New(MiddlewareWhen(func() (bool, error) { return condition, nil }, m))
	mux := http.NewServeMux()
	mux.Handle("/iftrue", appiftrue)
	mux.Handle("/iffalse", appiffalse)
	mux.Handle("/when", appwhen)
	server := httptest.NewServer(mux)
	defer server.Close()
	resp, err := http.Get(server.URL + "/iftrue")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	if resp.Header.Get("Condition") != "true" {
		t.Error(resp.Header.Get("Condition"))
	}
	resp, err = http.Get(server.URL + "/iffalse")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	if resp.Header.Get("Condition") != "" {
		t.Error(resp.Header.Get("Condition"))
	}
	resp, err = http.Get(server.URL + "/when")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	if resp.Header.Get("Condition") != "true" {
		t.Error(resp.Header.Get("Condition"))
	}
	condition = false
	resp, err = http.Get(server.URL + "/when")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	if resp.Header.Get("Condition") != "" {
		t.Error(resp.Header.Get("Condition"))
	}
}
