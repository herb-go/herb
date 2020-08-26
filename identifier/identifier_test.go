package identifier

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herb-go/herb/middleware"
)

var testIDentifier = IDFunc(func(r *http.Request) (string, error) {
	return r.Header.Get("id"), nil
})
var action = func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func TestIDentifier(t *testing.T) {
	appNil := middleware.New()
	appNil.Use(NewLoggedInFilter(testIDentifier, nil).ServeMiddleware).HandleFunc(action)
	appNotFound := middleware.New()
	appNotFound.Use(NewLoggedInFilter(testIDentifier, http.NotFoundHandler()).ServeMiddleware).HandleFunc(action)
	mux := &http.ServeMux{}
	mux.Handle("/nil", appNil)
	mux.Handle("/404", appNotFound)
	server := httptest.NewServer(mux)
	defer server.Close()
	req, err := http.NewRequest("GET", server.URL+"/nil", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Fatal()
	}
	resp.Body.Close()
	req, err = http.NewRequest("GET", server.URL+"/404", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 404 {
		t.Fatal(resp.StatusCode)
	}
	resp.Body.Close()
	req, err = http.NewRequest("GET", server.URL+"/404", nil)
	req.Header.Set("id", "test")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatal(resp.StatusCode)
	}
	resp.Body.Close()

}
