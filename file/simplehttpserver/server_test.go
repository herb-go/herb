package simplehttpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	var resp *http.Response
	var err error
	mux := http.NewServeMux()
	mux.HandleFunc("/test.html", ServeFile("./testdata/test.html"))
	mux.HandleFunc("/notexist.html", ServeFile("./testdata/notexist.html"))
	mux.Handle("/noindex/", http.StripPrefix("/noindex", http.HandlerFunc(ServeFolder("./testdata/noindex"))))
	mux.Handle("/index/", http.StripPrefix("/index", http.HandlerFunc(ServeFolder("./testdata/index"))))
	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err = http.DefaultClient.Get(server.URL + "/test.html")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}

	resp, err = http.DefaultClient.Get(server.URL + "/notexist.html")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Error(resp.StatusCode)
	}

	resp, err = http.DefaultClient.Get(server.URL + "/noindex/")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 403 {
		t.Error(resp.StatusCode)
	}

	resp, err = http.DefaultClient.Get(server.URL + "/noindex/noindex.html")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}

	resp, err = http.DefaultClient.Get(server.URL + "/index/")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}

	resp, err = http.DefaultClient.Get(server.URL + "/index/index.html")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	resp, err = http.DefaultClient.Get(server.URL + "/index/notexist.html")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Error(resp.StatusCode)
	}
}
