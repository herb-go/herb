package httpuser

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRedirector(t *testing.T) {

	redirector := NewRedirector("/target", "test", func(w http.ResponseWriter, r *http.Request) bool {
		return r.Header.Get("condition") == ""
	})
	var successAction = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}
	var targetAction = func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("condition") != "pass" {
			w.WriteHeader(422)
			w.Write([]byte("condition"))
		}
		url := redirector.MustClearSource(w, r)
		w.Write([]byte(url))
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirector.Middleware()(w, r, successAction)
	}))
	mux.HandleFunc("/target", targetAction)
	server := httptest.NewServer(mux)
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	client := http.Client{
		Jar: jar,
	}
	req, err := http.NewRequest("GET", server.URL+"/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.Request.URL.String() != server.URL+"/target" {
		t.Fatal(resp.Request.URL.String())
	}
	req, err = http.NewRequest("GET", server.URL+"/target", nil)
	req.Header.Set("condition", "pass")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(content) != "/test" {
		t.Fatal(string(content))
	}
	time.Sleep(100 * time.Millisecond)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	req, err = http.NewRequest("GET", server.URL+"/test", nil)
	req.Header.Set("condition", "pass")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

}
