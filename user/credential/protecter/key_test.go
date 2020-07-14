package protecter

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	id, err := DefaultKey.IdentifyRequest(r)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(id))
})

func TestForbiddenKey(t *testing.T) {
	s := httptest.NewServer(DefaultKey.ProtectWith(ForbiddenProtecter, testHandler))
	defer s.Close()
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 403 {
		t.Fatal(resp)
	}
}

func TestSuccessKey(t *testing.T) {
	s := httptest.NewServer(DefaultKey.ProtectWith(NotWorkingProtecter, testHandler))
	defer s.Close()
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
	if string(data) != "notworking" {
		t.Fatal()
	}
}

func TestNilKey(t *testing.T) {
	s := httptest.NewServer(DefaultKey.ProtectWith(nil, testHandler))
	defer s.Close()
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 403 {
		t.Fatal(resp)
	}
}
