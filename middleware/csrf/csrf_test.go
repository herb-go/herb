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

func newDefaultCsrf() *Csrf {
	c := New()
	config := Config{}
	config.Enabled = true
	err := config.ApplyTo(c)
	if err != nil {
		panic(err)
	}
	return c

}
func newTestCsrf() *Csrf {
	c := New()
	config := Config{
		CookieName: "herb-test-csrf-token",
		CookiePath: "/",
		HeaderName: "X-Csrf-Token",
		FormField:  "X-Csrf-Token",
		FailStatus: 400,
		Enabled:    true,
		FailHeader: defaultFailHeader,
		FailValue:  defaultFailValue,
	}
	err := config.ApplyTo(c)
	if err != nil {
		panic(err)
	}

	return c
}
func TestHeader(t *testing.T) {
	Csrf := newDefaultCsrf()
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
	mux.HandleFunc("/input", func(w http.ResponseWriter, r *http.Request) {
		output, err := Csrf.CsrfInput(w, r)
		if err != nil {
			t.Fatal(err)
		}
		w.Write([]byte(output))
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
	InputRequestWithToken, err := http.NewRequest("GET", s.URL+"/input", nil)
	token = jar.Cookies(InputRequestWithToken.URL)[0].Value
	InputRequestWithToken.Header.Set(Csrf.HeaderName, token)
	rep, err = c.Do(InputRequestWithToken)
	if err != nil {
		t.Fatal(err)
	}
	body, err = ioutil.ReadAll(rep.Body)
	if err != nil {
		t.Fatal(err)
	}
	if rep.StatusCode != 200 || !strings.Contains(string(body), token) {
		t.Fatal("Csrf input fail", string(body))
	}
	Csrf.Enabled = false
	HeaderRequest, err = http.NewRequest("GET", s.URL+"/header", nil)
	if err != nil {
		t.Fatal(err)
	}
	rep, err = c.Do(HeaderRequest)
	if err != nil {
		t.Fatal(err)
	}
	body, err = ioutil.ReadAll(rep.Body)
	if err != nil {
		t.Fatal(err)
	}
	rep.Body.Close()
	if rep.StatusCode != 200 || string(body) != successMsg {
		t.Errorf("Csrf block fail")
	}

}

func TestForm(t *testing.T) {
	Csrf := newTestCsrf()
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
