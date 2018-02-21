package basicauth

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/herb-go/herb/middleware"
)

func TestSingleUser(t *testing.T) {
	var content []byte
	var usernameAction = func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(GetUsername(r)))
		if err != nil {
			panic(err)
		}
	}
	var idAction = func(w http.ResponseWriter, r *http.Request) {
		id, err := Username.IdentifyRequest(r)
		if err != nil {
			panic(err)
		}
		_, err = w.Write([]byte(id))
		if err != nil {
			panic(err)
		}

	}
	var noRealmConfig = &SingleUser{
		Realm:    "",
		Username: "",
		Password: "",
	}
	var username = "testusername"
	var password = "testpassword"
	var usernameWrong = "testusernamewrong"
	var passwordWrong = "testpasswordwrong"
	var RealmConfig = &SingleUser{
		Realm:    "testrealm",
		Username: username,
		Password: password,
	}
	var NoRealmApp = middleware.New(Middleware(noRealmConfig)).HandleFunc(usernameAction)
	var app = middleware.New(Middleware(RealmConfig))
	var mux = http.NewServeMux()

	var realmmux = http.NewServeMux()
	realmmux.HandleFunc("/username", usernameAction)
	realmmux.HandleFunc("/id", idAction)
	app.Handle(realmmux)
	mux.Handle("/norealm/", NoRealmApp)
	mux.Handle("/realm/", http.StripPrefix("/realm", app))

	server := httptest.NewServer(mux)
	defer server.Close()
	req, err := http.NewRequest("GET", server.URL+"/norealm/id", nil)
	req.SetBasicAuth("", "")
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Error(resp.StatusCode)
	}
	time.Sleep(1 * time.Millisecond)
	req, err = http.NewRequest("GET", server.URL+"/realm/id", nil)
	if err != nil {
		panic(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 401 {
		t.Error(resp.StatusCode)
	}
	time.Sleep(1 * time.Millisecond)
	req, err = http.NewRequest("GET", server.URL+"/realm/id", nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(username, passwordWrong)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 401 {
		t.Error(resp.StatusCode)
	}
	time.Sleep(1 * time.Millisecond)
	req, err = http.NewRequest("GET", server.URL+"/realm/id", nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(usernameWrong, password)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 401 {
		t.Error(resp.StatusCode)
	}
	time.Sleep(1 * time.Millisecond)
	req, err = http.NewRequest("GET", server.URL+"/realm/id", nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(username, password)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	if string(content) != username {
		t.Error(string(content))
	}
	time.Sleep(1 * time.Millisecond)
	req, err = http.NewRequest("GET", server.URL+"/realm/username", nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(username, password)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	if string(content) != username {
		t.Error(string(content))
	}
	time.Sleep(1 * time.Millisecond)

}
