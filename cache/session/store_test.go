package session

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMiddleware(t *testing.T) {
	s := getStore(time.Hour)
	defer s.Close()
	var testHeaderName = "testheader"
	s.Mode = StoreModeHeader
	s.CookieName = testHeaderName
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
		_, err := s.GetRequestSession(r)
		if err != ErrRequestTokenNotFound {
			t.Fatal(err)
		}
		w.Write([]byte("ok"))
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
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(content) == "" {
		t.Fatal(content)
	}
	ErrRequest, err := http.NewRequest("POST", hs.URL+"/err", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Do(ErrRequest)
	if err != nil {
		t.Fatal(err)
	}
	NormalRequest, err := http.NewRequest("POST", hs.URL+"/normal", nil)
	NormalRequest.Header.Set(testHeaderName, "test")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = c.Do(NormalRequest)
	if err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(content) != "test" {
		t.Fatal(content)
	}
}

func TestCookieMiddleware(t *testing.T) {
	s := getStore(3600)
	defer s.Close()
	var testCookieName = "testcookie"
	s.Mode = StoreModeCookie
	s.CookieName = testCookieName
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
		_, err := s.GetRequestSession(r)
		if err != ErrRequestTokenNotFound {
			t.Fatal(err)
		}
		w.Write([]byte("ok"))
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
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(content) == "" {
		t.Fatal(content)
	}
	ErrRequest, err := http.NewRequest("POST", hs.URL+"/err", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Do(ErrRequest)
	if err != nil {
		t.Fatal(err)
	}
	NormalRequest, err := http.NewRequest("POST", hs.URL+"/normal", nil)
	NormalRequest.Header.Set("Cookie", testCookieName+"=test")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = c.Do(NormalRequest)
	if err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(content) != "test" {
		t.Fatal(content)
	}
}

func TestAutoGenerate(t *testing.T) {
	s := getStore(3600)
	defer s.Close()
	var testHeaderName = "testheader"
	s.Mode = StoreModeHeader
	s.CookieName = testHeaderName
	s.AutoGenerate = true
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
	mux.HandleFunc("/auto", func(w http.ResponseWriter, r *http.Request) {
		s.InstallMiddleware()(w, r, Action)
	})
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
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(content) == "" {
		t.Fatal(content)
	}
}
func TestNotFoundError(t *testing.T) {
	var s = NewTestCacheStore(time.Hour)
	session, err := s.GenerateSession("test")
	if err != nil {
		t.Fatal(err)
	}
	result := ""
	err = session.Get("testval", &result)
	if err != ErrDataNotFound {
		t.Fatal(err)
	}
	if !s.IsNotFoundError(ErrDataNotFound) {
		t.Fatal(err)
	}
}

func TestMaxLifeTime(t *testing.T) {
	var s = NewTestCacheStore(time.Hour)
	defer s.Close()
	s.TokenMaxLifetime = 5 * time.Second
	s.TokenLifetime = 3 * time.Second
	s.UpdateActiveInterval = time.Millisecond
	session, err := s.GenerateSession("test")
	if err != nil {
		t.Fatal(err)
	}
	err = session.Set("testkey", "testvalue")
	if err != nil {
		t.Fatal(err)
	}
	err = session.Save()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(200 * time.Millisecond)
	for i := 0; i < 4; i++ {
		time.Sleep(900 * time.Millisecond)

		go func() {
			session, err = s.GenerateSession("test")
			if err != nil {
				t.Fatal(err)
			}
			err = s.LoadSession(session)
			if err != nil {
				t.Fatal(err)
			}
			err = session.Save()
			if err != nil {
				t.Fatal(err)
			}
		}()
	}

	time.Sleep(1900 * time.Millisecond)
	session, err = s.GenerateSession("test")
	if err != nil {
		t.Fatal(err)
	}
	err = s.LoadSession(session)
	if err != ErrDataNotFound {
		t.Fatal(err)
	}
	err = session.Save()
	if err != nil {
		t.Fatal(err)
	}
}
