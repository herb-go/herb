package csrf

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var successMsg = "ok"

func TestHeader(t *testing.T) {
	Csrf := New()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Csrf.ServeSetCsrfTokenMiddleware(w, r, func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(successMsg))
		})
	})
	mux.HandleFunc("/header", func(w http.ResponseWriter, r *http.Request) {
		Csrf.ServeVerifyHeaderMiddleware(w, r, func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(successMsg))
		})
	})
	s := httptest.NewServer(mux)
	defer s.Close()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	c := &http.Client{
		Jar: jar,
	}
	HeaderRequest, err := http.NewRequest("GET", s.URL+"/header", nil)
	if err != nil {
		t.Fatal(err)
	}
	rep, err := c.Do(HeaderRequest)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		t.Fatal(err)
	}
	if rep.StatusCode != 400 || string(body) == successMsg {
		t.Errorf("Csrf block fail")
	}

	SetCrsfRequest, err := http.NewRequest("GET", s.URL+"/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rep, err = c.Do(SetCrsfRequest)
	if err != nil {
		t.Fatal(err)
	}
	body, err = ioutil.ReadAll(rep.Body)
	if err != nil {
		t.Fatal(err)
	}
	if rep.StatusCode != 200 || string(body) != successMsg {
		t.Errorf("Set csrf token fail")
	}
	HeaderRequestWithToken, err := http.NewRequest("GET", s.URL+"/header", nil)
	token := jar.Cookies(HeaderRequestWithToken.URL)[0].Value
	HeaderRequestWithToken.Header.Set(Csrf.HeaderName, token)
	rep, err = c.Do(HeaderRequestWithToken)
	if err != nil {
		t.Fatal(err)
	}
	body, err = ioutil.ReadAll(rep.Body)
	if err != nil {
		t.Fatal(err)
	}
	if rep.StatusCode != 200 || string(body) != successMsg {
		t.Errorf("Csrf block fail")
	}
}

func TestForm(t *testing.T) {
	Csrf := New()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Csrf.ServeSetCsrfTokenMiddleware(w, r, func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(successMsg))
		})
	})
	mux.HandleFunc("/form", func(w http.ResponseWriter, r *http.Request) {
		Csrf.ServeVerifyFormMiddleware(w, r, func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(successMsg))
		})
	})
	s := httptest.NewServer(mux)
	defer s.Close()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	c := &http.Client{
		Jar: jar,
	}
	FormRequest, err := http.NewRequest("POST", s.URL+"/form", nil)
	if err != nil {
		t.Fatal(err)
	}
	rep, err := c.Do(FormRequest)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		t.Fatal(err)
	}
	if rep.StatusCode != 400 || string(body) == successMsg {
		t.Errorf("Csrf block fail")
	}

	SetCrsfRequest, err := http.NewRequest("GET", s.URL+"/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rep, err = c.Do(SetCrsfRequest)
	if err != nil {
		t.Fatal(err)
	}
	body, err = ioutil.ReadAll(rep.Body)
	if err != nil {
		t.Fatal(err)
	}
	if rep.StatusCode != 200 || string(body) != successMsg {
		t.Errorf("Set csrf token fail")
	}

	token := jar.Cookies(FormRequest.URL)[0].Value
	form := url.Values{}
	form.Set(Csrf.FormField, token)
	FormRequestWithToken, err := http.NewRequest("POST", s.URL+"/form", strings.NewReader(form.Encode()))
	FormRequestWithToken.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	rep, err = c.Do(FormRequestWithToken)
	if err != nil {
		t.Fatal(err)
	}
	body, err = ioutil.ReadAll(rep.Body)
	if err != nil {
		t.Fatal(err)
	}
	if rep.StatusCode != 200 || string(body) != successMsg {
		t.Errorf("Csrf block fail")
	}
}
