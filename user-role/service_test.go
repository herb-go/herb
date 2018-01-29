package role

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herb-go/herb/middleware"
)

type identifier struct {
	Headername string
}

func (i *identifier) IdentifyRequest(r *http.Request) (string, error) {
	return r.Header.Get(i.Headername), nil
}

type roleservice struct {
	Prefix string
}

func (s *roleservice) Roles(uid string) (*Roles, error) {
	return NewRoles(s.Prefix + uid), nil
}

func TestServuce(t *testing.T) {
	var testRule = "test"
	var testHeader = "testheader"
	var service = Service{
		RoleService: &roleservice{},
		Identifier: &identifier{
			Headername: testHeader,
		},
	}
	var app = middleware.New()
	app.
		Use(service.RolesAuthorizeOrForbiddenMiddleware(testRule)).
		HandleFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
	var s = httptest.NewServer(app)
	defer s.Close()
	res, err := http.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusForbidden {
		t.Error(res.StatusCode)
	}
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add(testHeader, testRule)
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Error(res.StatusCode)
	}
}
