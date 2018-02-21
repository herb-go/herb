package misc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herb-go/herb/middleware"
)

func TestNoCache(t *testing.T) {
	var successAction = func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		if err != nil {
			panic(err)
		}
	}
	app := middleware.New()
	app.Use(NoCache).HandleFunc(successAction)
	server := httptest.NewServer(app)
	defer server.Close()
	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.Header.Get("Expires") != "0" {
		t.Error(resp.Header.Get("Expires"))
	}
	if resp.Header.Get("Pragma") != "no-cache" {
		t.Error(resp.Header.Get("Pragma"))
	}
	if resp.Header.Get("Cache-Control") != "no-cache, no-store, must-revalidate" {
		t.Error(resp.Header.Get("Cache-Control"))
	}
}
