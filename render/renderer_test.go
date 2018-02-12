package render

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRenderer(t *testing.T) {
	renderer := New(nil, "")
	mux := http.NewServeMux()
	mux.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		renderer.MustJSON(w, "ok", 200)
	})
	mux.HandleFunc("/writejson", func(w http.ResponseWriter, r *http.Request) {
		renderer.MustWriteJSON(w, []byte(`"ok"`), 200)
	})
	mux.HandleFunc("/html", func(w http.ResponseWriter, r *http.Request) {
		renderer.MustWriteHTML(w, []byte("ok"), 200)
	})
	mux.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) {
		renderer.MustError(w, 404)
	})
	mux.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		renderer.MustHTMLFile(w, "testdata/test.html", 200)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	resp, err := http.DefaultClient.Get(server.URL + "/json")
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(content) != `"ok"` {
		t.Error(content)
	}
	if resp.Header.Get(ContentType) != ContentJSON {
		t.Error(resp.Header.Get(ContentType))
	}

	resp, err = http.DefaultClient.Get(server.URL + "/writejson")
	if err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(content) != `"ok"` {
		t.Error(content)
	}
	if resp.Header.Get(ContentType) != ContentJSON {
		t.Error(resp.Header.Get(ContentType))
	}

	resp, err = http.DefaultClient.Get(server.URL + "/html")
	if err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(content) != `ok` {
		t.Error(content)
	}
	if resp.Header.Get(ContentType) != ContentHTML {
		t.Error(resp.Header.Get(ContentType))
	}

	resp, err = http.DefaultClient.Get(server.URL + "/404")
	if err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(content) != http.StatusText(404) {
		t.Error(content)
	}
	if resp.StatusCode != 404 {
		t.Error(resp.StatusCode)
	}
	if resp.Header.Get(ContentType) != ContentText {
		t.Error(resp.Header.Get(ContentType))
	}

	resp, err = http.DefaultClient.Get(server.URL + "/file")
	if err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(content) != `ok` {
		t.Error(content)
	}
	if resp.Header.Get(ContentType) != ContentHTML {
		t.Error(resp.Header.Get(ContentType))
	}
}
