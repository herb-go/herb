package user

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

type testIdentifier struct {
	cookieName string
}

func (i *testIdentifier) IdentifyRequest(r *http.Request) (string, error) {
	c, err := r.Cookie(i.cookieName)
	if err == http.ErrNoCookie {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return c.Value, nil
}
func (i *testIdentifier) Logout(w http.ResponseWriter, r *http.Request) error {
	c := &http.Cookie{
		Name:   i.cookieName,
		MaxAge: 0,
		Value:  "",
		Path:   "/",
	}
	http.SetCookie(w, c)
	return nil
}
func (i *testIdentifier) Login(w http.ResponseWriter, r *http.Request, id string) error {
	c := &http.Cookie{
		Name:   i.cookieName,
		MaxAge: 3600,
		Value:  id,
		Path:   "/",
	}
	http.SetCookie(w, c)
	return nil
}
func newTestIdentifier(cookieName string) *testIdentifier {
	return &testIdentifier{
		cookieName: cookieName,
	}
}

func TestLoginRequired(t *testing.T) {
	Identifier := newTestIdentifier("test")
	loginRequired := LoginRequiredMiddleware(Identifier, nil)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		loginRequired(w, r, func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		err := Identifier.Login(w, r, r.Header.Get("loginid"))
		if err != nil {
			panic(err)
		}
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		LogoutMiddleware(Identifier)(w, r, func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
	})
	mux.HandleFunc("/id2", func(w http.ResponseWriter, r *http.Request) {
		MiddlewareForbiddenExceptForUsers(Identifier, []string{"id2"})(w, r, func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	req, err := http.NewRequest("get", server.URL+"/", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Error(resp.StatusCode)
	}
	resp.Body.Close()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	client := http.Client{
		Jar: jar,
	}
	loginreq, err := http.NewRequest("get", server.URL+"/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	loginreq.Header.Add("loginid", "testloginid")
	resp, err = client.Do(loginreq)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	resp.Body.Close()
	req, err = http.NewRequest("get", server.URL+"/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req, err = http.NewRequest("get", server.URL+"/", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	resp.Body.Close()

	logoutreq, err := http.NewRequest("get", server.URL+"/logout", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(logoutreq)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	resp.Body.Close()
	req, err = http.NewRequest("get", server.URL+"/", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Error(resp.StatusCode)
	}
	resp.Body.Close()
	reqid2, err := http.NewRequest("get", server.URL+"/id2", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(reqid2)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 403 {
		t.Error(resp.StatusCode)
	}
	resp.Body.Close()
	loginreq, err = http.NewRequest("get", server.URL+"/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	loginreq.Header.Add("loginid", "testloginid")
	resp, err = client.Do(loginreq)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	resp.Body.Close()
	reqid2, err = http.NewRequest("get", server.URL+"/id2", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(reqid2)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 403 {
		t.Error(resp.StatusCode)
	}
	resp.Body.Close()
	loginreq, err = http.NewRequest("get", server.URL+"/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	loginreq.Header.Add("loginid", "id2")
	resp, err = client.Do(loginreq)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	resp.Body.Close()
	reqid2, err = http.NewRequest("get", server.URL+"/id2", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(reqid2)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	resp.Body.Close()
}

func TestLoginRequiredCustomFailAction(t *testing.T) {
	failAction := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(400), 400)
	}
	Identifier := newTestIdentifier("test")
	loginRequired := LoginRequiredMiddleware(Identifier, failAction)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		loginRequired(w, r, func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		err := Identifier.Login(w, r, r.Header.Get("loginid"))
		if err != nil {
			panic(err)
		}
		w.Write([]byte("ok"))
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	req, err := http.NewRequest("get", server.URL+"/", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Error(resp.StatusCode)
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	client := http.Client{
		Jar: jar,
	}
	loginreq, err := http.NewRequest("get", server.URL+"/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	loginreq.Header.Add("loginid", "testloginid")
	resp, err = client.Do(loginreq)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	req, err = http.NewRequest("get", server.URL+"/", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}

}

func TestLoginLoginRedirect(t *testing.T) {
	Identifier := newTestIdentifier("test")
	loginRedirector := NewLoginRedirector("/redirect", "redirect")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		loginRedirector.Middleware(Identifier)(w, r, func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		err := Identifier.Login(w, r, r.Header.Get("loginid"))
		if err != nil {
			panic(err)
		}
		url := loginRedirector.MustClearSource(w, r)
		w.Write([]byte(url))
	})
	mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(422), 422)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	client := http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("get", server.URL+"/", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 422 {
		t.Error(resp.StatusCode)
	}
	loginreq, err := http.NewRequest("get", server.URL+"/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	loginreq.Header.Add("loginid", "testloginid")
	resp, err = client.Do(loginreq)
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	if string(content) != "/" {
		t.Error(string(content))
	}
	req, err = http.NewRequest("get", server.URL+"/", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	loginreq, err = http.NewRequest("get", server.URL+"/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	loginreq.Header.Add("loginid", "testloginid")
	resp, err = client.Do(loginreq)
	if err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	if string(content) != "" {
		t.Error(string(content))
	}

}
