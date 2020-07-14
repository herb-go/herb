package protecter

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herb-go/herb/user/credential"
)

var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	id, err := DefaultKey.IdentifyRequest(r)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(id))
})
var testIDHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(LoadID(r)))
})
var credentialerAppID = CredentialerFunc(func(r *http.Request) credential.Credential {
	return credential.New().WithType(credential.TypeAppID).WithData([]byte(r.Header.Get("appid")))
})

var credentialerToken = CredentialerFunc(func(r *http.Request) credential.Credential {
	return credential.New().WithType(credential.TypeToken).WithData([]byte(r.Header.Get("token")))
})

var notfound = http.NotFoundHandler()

var testProtecter = New().
	WithOnFail(notfound).
	WithCredentialers(credentialerAppID, credentialerToken).
	WithVerifier(
		credential.VerifierFunc(func(c *credential.Collection) (string, error) {
			if string(c.Get(credential.TypeAppID)) == "testappid" && string(c.Get(credential.TypeToken)) == "testtoken" {
				return "testappid", nil
			}
			return "", nil
		},
			credential.TypeAppID,
			credential.TypeToken,
		),
	)

func TestForbidden(t *testing.T) {
	s := httptest.NewServer(ProtectWith(ForbiddenProtecter, testHandler))
	defer s.Close()
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 403 {
		t.Fatal(resp)
	}
}

func TestSuccess(t *testing.T) {
	s := httptest.NewServer(ProtectWith(NotWorkingProtecter, testHandler))
	defer s.Close()
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)

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

func TestNil(t *testing.T) {
	s := httptest.NewServer(ProtectWith(nil, testHandler))
	defer s.Close()
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 403 {
		t.Fatal(resp)
	}
}

func TestVerifyFail(t *testing.T) {
	s := httptest.NewServer(ProtectWith(testProtecter, testHandler))
	defer s.Close()
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Fatal(resp)
	}

}

func TestVerifySuccess(t *testing.T) {
	s := httptest.NewServer(ProtectWith(testProtecter, testIDHandler))
	defer s.Close()
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("appid", "testappid")
	req.Header.Add("token", "testtoken")
	resp, err := http.DefaultClient.Do(req)
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
	if string(data) != "testappid" {
		t.Fatal(string(data))
	}
}

func TestMiddlewareFail(t *testing.T) {
	m := ProtectMiddleware(nil)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m(w, r, testHandler)
	}))
	defer s.Close()
	resp, err := http.Get(s.URL)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 403 {
		t.Fatal(resp)
	}
}

func TestMiddlewareSuccess(t *testing.T) {
	m := ProtectMiddleware(NotWorkingProtecter)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m(w, r, testHandler)
	}))
	defer s.Close()
	resp, err := http.Get(s.URL)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
}
