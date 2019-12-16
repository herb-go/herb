package misc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herb-go/herb/middleware"
)

func TestMethod(t *testing.T) {
	var successAction = func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		if err != nil {
			panic(err)
		}
	}
	app := middleware.New()
	app.Use(MethodMiddleware("GET")).HandleFunc(successAction)
	server := httptest.NewServer(app)
	defer server.Close()
	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
	resp.Body.Close()
	resp, err = http.Post(server.URL+"/", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 405 {
		t.Fatal(resp)
	}

}
