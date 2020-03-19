package session

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestCookieWriter(t *testing.T) {
	s := getStore(time.Hour)
	defer s.Close()
	var testCookieName = "testcookie"
	s.Mode = StoreModeCookie
	s.CookieName = testCookieName
	s.UpdateActiveInterval = 1 * time.Microsecond
	s.TokenLifetime = 10 * time.Second
	s.TokenMaxLifetime = 100 * time.Second
	var mux = http.NewServeMux()
	var Action = func(w http.ResponseWriter, r *http.Request) {
		session, err := s.GetRequestSession(r)
		if err != nil {
			panic(err)
		}
		token, err := session.Token()
		if err != nil {
			panic(err)
		}
		w.Write([]byte(token))
	}
	var ActionErr = func(w http.ResponseWriter, r *http.Request) {
		s, err := s.GetRequestSession(r)
		if err != ErrRequestTokenNotFound {
			t.Fatal(err)
		}
		w.Write([]byte(strconv.FormatInt(s.ExpiredAt, 10)))
	}
	mux.HandleFunc("/auto", func(w http.ResponseWriter, r *http.Request) {
		s.InstallMiddleware()(w, r, func(w http.ResponseWriter, r *http.Request) {
			s.AutoGenerateMiddleware()(w, r, Action)
		})
	})
	mux.HandleFunc("/normal", func(w http.ResponseWriter, r *http.Request) {
		s.InstallMiddleware()(w, r, Action)
	})
	mux.HandleFunc("/err", ActionErr)
	hs := httptest.NewServer(mux)
	defer hs.Close()
	c := &http.Client{}
	AutoRequest, err := http.NewRequest("POST", hs.URL+"/auto", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Do(AutoRequest)
	if err != nil {
		t.Fatal(err)
	}
	content1, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	time.Sleep(100 * time.Millisecond)
	resp, err = c.Do(AutoRequest)
	if err != nil {
		t.Fatal(err)
	}
	content2, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(content1) == string(content2) {
		t.Fatal(string(content1), string(content2))
	}

}
