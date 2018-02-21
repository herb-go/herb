package misc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herb-go/herb/middleware"
)

func TestPoweredBy(t *testing.T) {
	var successAction = func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		if err != nil {
			panic(err)
		}
	}
	app := middleware.New()
	app.Use(PoweredBy("herb")).HandleFunc(successAction)
	server := httptest.NewServer(app)
	defer server.Close()
	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.Header.Get("Powered-By") != "herb" {
		t.Error(resp.Header.Get("Powered-By"))
	}
}
