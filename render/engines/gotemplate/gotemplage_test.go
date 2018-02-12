package gotemplate

import (
	"encoding/base64"
	"testing"

	"github.com/herb-go/herb/render"
)

func TestTemplate(t *testing.T) {
	engine := Engine
	var b64 = func(data string) (string, error) {
		d := base64.RawStdEncoding.EncodeToString([]byte(data))
		return d, nil
	}
	engine.SetViewRoot("./testdata")
	engine.FuncMap["b64"] = b64

	view, err := engine.Compile("test.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	data := new(render.Data)
	data.Set("data", "testdata")
	data.Set("raw", "<html/>")

	output, err := view.Execute(data)
	if err != nil {
		t.Fatal(err)
	}
	result := base64.RawStdEncoding.EncodeToString([]byte("testdata"))
	if string(output) != result {
		t.Error(output)
	}
	viewraw, err := engine.Compile("raw.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	outputraw, err := viewraw.Execute(data)
	if err != nil {
		t.Fatal(err)
	}
	if string(outputraw) != "<html/>" {
		t.Error(outputraw)
	}
}
