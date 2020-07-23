package shouldvalidated

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herb-go/herb/middleware"
	"github.com/herb-go/herb/middleware/httpinfo"
)

func testValidator(data []byte) (bool, error) {
	return string(data) == "ok", nil
}

var testField = httpinfo.FieldFunc(func(r *http.Request) ([]byte, bool, error) {
	if r.Header.Get("test") == "false" {
		return nil, false, nil
	}
	return []byte(r.Header.Get("test")), true, nil
})

var action = func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func TestMiddleware(t *testing.T) {
	app := middleware.New()
	app.Use(NewNotFoundMiddleware(testField, testValidator)).HandleFunc(action)
	s := httptest.NewServer(app)
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
	req, err = http.NewRequest("GET", s.URL, nil)
	req.Header.Add("test", "false")
	if err != nil {
		panic(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Fatal(resp)
	}
	req, err = http.NewRequest("GET", s.URL, nil)
	req.Header.Add("test", "ok")
	if err != nil {
		panic(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}

}
