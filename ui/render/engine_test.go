package render

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testEngine struct {
	Root string
}

func (e *testEngine) SetViewRoot(path string) {
	e.Root = path
}
func (e *testEngine) Compile(config *ViewConfig) (CompiledView, error) {
	return &testView{
		Files: config.Files,
	}, nil
}
func (e *testEngine) RegisterFunc(name string, fn interface{}) error {
	return ErrRegisterFuncNotSupported
}

type testView struct {
	Files []string
}

func (v *testView) Execute(data interface{}) ([]byte, error) {
	d := data.(Data)
	bs := d.Get("data").([]byte)
	return bs, nil
}

func TestEngine(t *testing.T) {
	engine := &testEngine{}
	render := New()
	oc := NewOptionCommon()
	oc.Engine = engine
	err := render.Init(oc)
	if err != nil {
		panic(err)
	}
	if render.Engine() != engine {
		t.Error(render.Engine())
	}
	var b64 = func(data string) string {
		d := base64.RawStdEncoding.EncodeToString([]byte(data))
		return d
	}
	err = engine.RegisterFunc("b64", b64)
	if err != ErrRegisterFuncNotSupported {
		t.Error(err)
	}
	option := ViewsOptionCommon{}
	bytes, err := ioutil.ReadFile("./testdata/testconfig.json")
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(bytes, &option)
	if err != nil {
		t.Error(err)
	}
	render.MustInitViews(option)
	ViewTest := render.GetView("test")
	ViewTestNew := render.NewView("testnew", NewViewConfig(
		"testnew.view",
	))
	ViewNotExist := render.GetView("testnotexist")

	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		data := Data{}
		data.Set("data", []byte("testdata"))
		ViewTest.MustRender(w, data)
	})
	mux.HandleFunc("/testnew", func(w http.ResponseWriter, r *http.Request) {
		data := Data{}
		data.Set("data", []byte("testnewdata"))
		_, err := ViewTestNew.Render(w, data)
		if err != nil {
			panic(err)
		}
	})
	mux.HandleFunc("/testnotexist", func(w http.ResponseWriter, r *http.Request) {
		data := Data{}
		data.Set("data", []byte("testnotexistdata"))
		_, err := ViewNotExist.Render(w, data)
		e, ok := err.(*ViewError)
		if ok && e.err == ErrViewNotExist {
			http.Error(w, err.Error(), 400)
			return
		}
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	resp, err := http.DefaultClient.Get(server.URL + "/test")
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(content) != "testdata" {
		t.Error(string(content))
	}

	resp, err = http.DefaultClient.Get(server.URL + "/testnew")
	if err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(content) != "testnewdata" {
		t.Error(string(content))
	}

	resp, err = http.DefaultClient.Get(server.URL + "/testnotexist")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 400 {
		t.Error(resp.StatusCode)
	}
}
