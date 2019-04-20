package captcha

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCaptcha(t *testing.T) {
	c := NewCatpcha()
	mux := http.NewServeMux()
	mux.HandleFunc("/captcha", c.CaptchaAction("test"))
	mux.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		verifier := c.Verifier(r, "test")
		code := r.Header.Get("code")
		result, err := verifier(code)
		if err != nil {
			t.Fatal(err)
		}
		if result {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(422)
		}
		w.Write([]byte(code))
	})
	hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.SessionStore.HeaderMiddleware("token")(w, r, mux.ServeHTTP)
	}))
	defer hs.Close()
	ts, err := c.SessionStore.GenerateSession("testtoken")
	if err != nil {
		t.Fatal(err)
	}
	ts.Save()
	req, err := http.NewRequest("GET", hs.URL+"/captcha", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("token", "testtoken")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	result := map[string]interface{}{}
	err = json.Unmarshal(content, &result)
	if err != nil {
		t.Fatal(err)
	}
	if result["Code"] == "" {
		t.Fatal(result["Code"])
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	result2 := map[string]interface{}{}
	err = json.Unmarshal(content, &result2)
	if err != nil {
		t.Fatal(err)
	}
	if result2["Code"] != result["Code"] {
		t.Fatal(result2["Code"])
	}
	time.Sleep(10 * time.Millisecond)
	req, err = http.NewRequest("GET", hs.URL+"/captcha", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("token", "testtoken")
	req.Header.Set(HeaderReset, "true")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	result3 := map[string]interface{}{}
	err = json.Unmarshal(content, &result3)
	if err != nil {
		t.Fatal(err)
	}
	if result3["Code"] == result2["Code"] {
		t.Fatal(result3)
	}
	req, err = http.NewRequest("GET", hs.URL+"/verify", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("token", "testtoken")
	req.Header.Set("code", result3["Code"].(string))
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
	req, err = http.NewRequest("GET", hs.URL+"/verify", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("token", "testtoken")
	req.Header.Set("code", result2["Code"].(string))
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 422 {
		t.Fatal(resp)
	}
}

func TestDiasbledCaptcha(t *testing.T) {
	captcha := newEmptyCaptcha()
	c := &Config{
		Driver:  "testcaptcha",
		Enabled: false,
	}
	err := c.ApplyTo(captcha)
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/captcha", captcha.CaptchaAction("test"))
	mux.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		verifier := captcha.Verifier(r, "test")
		code := r.Header.Get("code")
		result, err := verifier(code)
		if err != nil {
			t.Fatal(err)
		}
		if result {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(422)
		}
		w.Write([]byte(code))
	})
	hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captcha.SessionStore.HeaderMiddleware("token")(w, r, mux.ServeHTTP)
	}))
	defer hs.Close()
	ts, err := captcha.SessionStore.GenerateSession("testtoken")
	if err != nil {
		t.Fatal(err)
	}
	ts.Save()
	req, err := http.NewRequest("GET", hs.URL+"/captcha", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("token", "testtoken")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "{}" {
		t.Fatal(string(content))

	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
	req, err = http.NewRequest("GET", hs.URL+"/verify", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("token", "testtoken")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}

}

func TestDiasbledSceneCaptcha(t *testing.T) {
	captcha := newEmptyCaptcha()
	c := &Config{
		Driver:         "testcaptcha",
		Enabled:        true,
		DisabledScenes: map[string]bool{"test": true},
	}
	err := c.ApplyTo(captcha)
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/captcha", captcha.CaptchaAction("test"))
	mux.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		verifier := captcha.Verifier(r, "test")
		code := r.Header.Get("code")
		result, err := verifier(code)
		if err != nil {
			t.Fatal(err)
		}
		if result {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(422)
		}
		w.Write([]byte(code))
	})
	hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captcha.SessionStore.HeaderMiddleware("token")(w, r, mux.ServeHTTP)
	}))
	defer hs.Close()
	ts, err := captcha.SessionStore.GenerateSession("testtoken")
	if err != nil {
		t.Fatal(err)
	}
	ts.Save()
	req, err := http.NewRequest("GET", hs.URL+"/verify", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("token", "testtoken")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
}

func TestIpWhiteListCaptcha(t *testing.T) {
	captcha := newEmptyCaptcha()
	c := &Config{
		Driver:        "testcaptcha",
		Enabled:       true,
		AddrWhiteList: []string{"127.0.0.1"},
	}
	err := c.ApplyTo(captcha)
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/captcha", captcha.CaptchaAction("test"))
	mux.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		r.RemoteAddr = "127.0.0.1:8000"
		verifier := captcha.Verifier(r, "test")
		code := r.Header.Get("code")
		result, err := verifier(code)
		if err != nil {
			t.Fatal(err)
		}
		if result {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(422)
		}
		w.Write([]byte(code))
	})
	hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captcha.SessionStore.HeaderMiddleware("token")(w, r, mux.ServeHTTP)
	}))
	defer hs.Close()
	ts, err := captcha.SessionStore.GenerateSession("testtoken")
	if err != nil {
		t.Fatal(err)
	}
	ts.Save()
	req, err := http.NewRequest("GET", hs.URL+"/verify", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("token", "testtoken")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
}
