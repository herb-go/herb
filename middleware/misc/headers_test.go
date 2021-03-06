package misc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herb-go/herb/middleware"
)

func TestHeaders(t *testing.T) {
	var successAction = func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		if err != nil {
			panic(err)
		}
	}
	app := middleware.New()
	m := Headers{
		"X-Powered-By": "Herbgo",
	}
	app.Use(m.ServeMiddleware).HandleFunc(successAction)
	server := httptest.NewServer(app)
	defer server.Close()
	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	h := resp.Header.Get("X-Powered-By")
	if h != "Herbgo" {
		t.Error(h)
	}
}
