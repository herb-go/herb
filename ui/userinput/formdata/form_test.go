package formdata

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	model "github.com/herb-go/herb/ui/userinput"
)

type testForm struct {
	Form
	Field1 string
	User   string
}

func (f *testForm) Validate() error {
	f.ValidateField(f.Field1 != "", "Field1", "Field1 required.")
	return nil
}
func (f *testForm) InitWithRequest(r *http.Request) error {
	f.User = r.Header.Get("user")
	return f.Form.InitWithRequest(r)
}

func newTestForm() *testForm {
	return &testForm{}
}
func TestForm(t *testing.T) {
	var testAction = func(w http.ResponseWriter, r *http.Request) {
		form := newTestForm()
		if MustValidateJSONRequest(r, form) {
			if form.HTTPRequest() != r {
				t.Fatal(r)
			}
			bytes, err := json.Marshal(form)
			if err != nil {
				panic(err)
			}
			_, err = w.Write(bytes)
			if err != nil {
				panic(err)
			}
		} else {
			MustRenderErrorsJSON(w, form)
		}
	}

	var req *http.Request
	var resp *http.Response
	var content []byte
	var err error
	server := httptest.NewServer(http.HandlerFunc(testAction))
	defer server.Close()
	req, err = http.NewRequest("POST", server.URL+"/", nil)
	if err != nil {
		t.Fatal(err)
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
	if resp.StatusCode != 422 {
		t.Error(resp.StatusCode)
	}
	var errdata = []model.FieldError{}
	err = json.Unmarshal(content, &errdata)
	if err != nil {
		t.Fatal(err)
	}
	if len(errdata) != 1 {
		t.Error(errdata)
	}
	time.Sleep(time.Millisecond)

	req, err = http.NewRequest("POST", server.URL+"/", bytes.NewBuffer([]byte(`non json data`)))
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Error(resp.StatusCode)
	}
	time.Sleep(time.Millisecond)

	form := newTestForm()
	form.Field1 = "Field1"
	formdata, err := json.Marshal(form)
	req, err = http.NewRequest("POST", server.URL+"/", bytes.NewBuffer(formdata))
	req.Header.Add("user", "testuser")
	if err != nil {
		t.Fatal(err)
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
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	var resultform = newTestForm()
	err = json.Unmarshal(content, resultform)
	if err != nil {
		t.Fatal(err)
	}
	if resultform.Field1 != form.Field1 {
		t.Error(resultform.Field1)
	}
	if resultform.User != "testuser" {
		t.Error(resultform.User)
	}
}
