package spilter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herb-go/herb/middleware"
	"github.com/herb-go/herb/middleware/router"
)

var action = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("test", router.GetParams(r).Get("test"))
	w.Header().Set("path", r.URL.Path)
	w.Write([]byte("ok"))
}

func TestSplitFirstFolderInto1(t *testing.T) {
	app := middleware.New()
	app.Use(SplitFirstFolderInto("test")).HandleFunc(action)
	s := httptest.NewServer(app)
	defer s.Close()
	resp, err := http.Get(s.URL + "/testpath")
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.Header.Get("test") != "testpath" {
		t.Fatal(resp.Header.Get("test"))
	}
	if resp.Header.Get("path") != "/" {
		t.Fatal(resp.Header.Get("path"))
	}
}

func TestSplitFirstFolderInto2(t *testing.T) {
	app := middleware.New()
	app.Use(SplitFirstFolderInto("test")).HandleFunc(action)
	s := httptest.NewServer(app)
	defer s.Close()
	resp, err := http.Get(s.URL + "/testpath/left")
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.Header.Get("test") != "testpath" {
		t.Fatal(resp.Header.Get("test"))
	}
	if resp.Header.Get("path") != "/left" {
		t.Fatal(resp.Header.Get("path"))
	}
}
func TestDropAfterFirst1(t *testing.T) {
	app := middleware.New()
	app.Use(DropAfterFirst("-")).HandleFunc(action)
	s := httptest.NewServer(app)
	defer s.Close()
	resp, err := http.Get(s.URL + "/testpath-left")
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	if resp.Header.Get("path") != "/testpath" {
		t.Fatal(resp.Header.Get("path"))
	}
}

func TestDropAfterFirst2(t *testing.T) {
	app := middleware.New()
	app.Use(DropAfterFirst("-")).HandleFunc(action)
	s := httptest.NewServer(app)
	defer s.Close()
	resp, err := http.Get(s.URL + "/testpath")
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	if resp.Header.Get("path") != "/testpath" {
		t.Fatal(resp.Header.Get("path"))
	}
}
