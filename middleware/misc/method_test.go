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

func TestMethos(t *testing.T) {
	var successAction = func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		if err != nil {
			panic(err)
		}
	}
	mux := http.ServeMux{}
	mux.HandleFunc("/GET", func(w http.ResponseWriter, r *http.Request) {
		MethodGET.ServeMiddleware(w, r, successAction)
	})
	mux.HandleFunc("/POST", func(w http.ResponseWriter, r *http.Request) {
		MethodPOST.ServeMiddleware(w, r, successAction)
	})
	mux.HandleFunc("/PUT", func(w http.ResponseWriter, r *http.Request) {
		MethodPUT.ServeMiddleware(w, r, successAction)
	})
	mux.HandleFunc("/DELETE", func(w http.ResponseWriter, r *http.Request) {
		MethodDELETE.ServeMiddleware(w, r, successAction)
	})
	mux.HandleFunc("/OPTIONS", func(w http.ResponseWriter, r *http.Request) {
		MethodOPTIONS.ServeMiddleware(w, r, successAction)
	})
	s := httptest.NewServer(&mux)
	defer s.Close()
	resp, err := http.DefaultClient.Get(s.URL + "/GET")
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
	resp, err = http.DefaultClient.Post(s.URL+"/GET", "", nil)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 405 {
		t.Fatal(resp)
	}
	resp, err = http.DefaultClient.Post(s.URL+"/POST", "", nil)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
	resp, err = http.DefaultClient.Get(s.URL + "/POST")
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 405 {
		t.Fatal(resp)
	}

	req, err := http.NewRequest("PUT", s.URL+"/PUT", nil)
	if err != nil {
		panic(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}

	resp, err = http.DefaultClient.Get(s.URL + "/PUT")
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 405 {
		t.Fatal(resp)
	}

	req, err = http.NewRequest("DELETE", s.URL+"/DELETE", nil)
	if err != nil {
		panic(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
	resp, err = http.DefaultClient.Get(s.URL + "/DELETE")
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 405 {
		t.Fatal(resp)
	}
	req, err = http.NewRequest("OPTIONS", s.URL+"/OPTIONS", nil)
	if err != nil {
		panic(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
	resp, err = http.DefaultClient.Get(s.URL + "/OPTIONS")
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 405 {
		t.Fatal(resp)
	}
}
