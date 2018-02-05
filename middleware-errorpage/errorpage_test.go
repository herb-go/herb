package errorpage

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/herb-go/herb/middleware"
)

func TestErrorPage(t *testing.T) {
	app := middleware.New()
	mux := http.NewServeMux()
	errorpage := New()
	errorpage.OnError(func(w http.ResponseWriter, r *http.Request, status int) {
		w.Write([]byte("error:" + strconv.Itoa(status)))
	})
	errorpage.OnStatus(404, func(w http.ResponseWriter, r *http.Request, status int) {
		w.Write([]byte(strconv.Itoa(status)))
	})
	errorpage.IgnoreStatus(403)
	mux.HandleFunc("/200", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		if err != nil {
			panic(err)
		}
	})
	mux.HandleFunc("/403", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(403), 403)
	})
	mux.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(404), 404)
	})
	mux.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(500), 500)
	})
	mux.Handle("/disable404", middleware.New(errorpage.MiddlewareDisable).HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(404), 404)
	}))
	mux.Handle("/disable500", middleware.New(errorpage.MiddlewareDisable).HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(500), 500)
	}))
	app.Use(errorpage.ServeMiddleware).Handle(mux)
	server := httptest.NewServer(app)
	defer server.Close()
	resp, err := http.Get(server.URL + "/404")
	if err != nil {
		panic(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		panic(err)
	}
	if string(content) != "404" {
		t.Error(string(content))
	}
	resp, err = http.Get(server.URL + "/500")
	if err != nil {
		panic(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		panic(err)
	}
	if string(content) != "error:500" {
		t.Error(string(content))
	}
	resp, err = http.Get(server.URL + "/200")
	if err != nil {
		panic(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		panic(err)
	}
	if string(content) != "ok" {
		t.Error(string(content))
	}
	resp, err = http.Get(server.URL + "/403")
	if err != nil {
		panic(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		panic(err)
	}
	if string(content) != http.StatusText(403)+"\n" {
		t.Error(string(content))
	}
	resp, err = http.Get(server.URL + "/disable404")
	if err != nil {
		panic(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		panic(err)
	}
	if string(content) != http.StatusText(404)+"\n" {
		t.Error(string(content))
	}
	resp, err = http.Get(server.URL + "/disable500")
	if err != nil {
		panic(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		panic(err)
	}
	if string(content) != http.StatusText(500)+"\n" {
		t.Error(string(content))
	}
}
