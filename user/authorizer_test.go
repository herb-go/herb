package user

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type testAuthorizer struct {
	name  string
	value string
}

func newTestAuthorizer(name string, value string) *testAuthorizer {
	return &testAuthorizer{
		name:  name,
		value: value,
	}
}
func (a *testAuthorizer) Authorize(r *http.Request) (bool, error) {
	return r.Header.Get(a.name) == a.value, nil
}

func successAction(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("ok"))
	if err != nil {
		panic(err)
	}
}

func TestAuthorizer(t *testing.T) {
	Authorizer := newTestAuthorizer("testname", "testvalue")
	failAction := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(400), 400)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AuthorizeMiddleware(Authorizer, failAction)(w, r, successAction)
	}))
	defer server.Close()
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Error(resp.StatusCode)
	}
	time.Sleep(1 * time.Millisecond)
	req.Header.Add("testname", "testvalue")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
}

func TestAuthorizerNilFailAction(t *testing.T) {
	Authorizer := newTestAuthorizer("testname", "testvalue")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AuthorizeMiddleware(Authorizer, nil)(w, r, successAction)
	}))
	defer server.Close()
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 403 {
		t.Error(resp.StatusCode)
	}
	time.Sleep(1 * time.Millisecond)
	req.Header.Add("testname", "testvalue")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
}
func TestAuthorizeOrForbiddenMiddleware(t *testing.T) {
	Authorizer := newTestAuthorizer("testname", "testvalue")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AuthorizeOrForbiddenMiddleware(Authorizer)(w, r, successAction)
	}))
	defer server.Close()
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 403 {
		t.Error(resp.StatusCode)
	}
	time.Sleep(1 * time.Millisecond)
	req.Header.Add("testname", "testvalue")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
}
