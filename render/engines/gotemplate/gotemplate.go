package gotemplate

import (
	"html/template"
	"path/filepath"

	"bytes"

	"github.com/herb-go/herb/render"
)

func unescaped(x string) interface{} { return template.HTML(x) }

func New() *GoTemplateEngine {
	e := GoTemplateEngine{
		FuncMap: template.FuncMap{
			"Raw": unescaped,
		},
	}
	return &e
}

type GoTemplateView template.Template

func (v *GoTemplateView) Execute(data interface{}) (string, error) {

	writer := bytes.NewBuffer([]byte{})
	t := template.Template(*v)
	err := t.Execute(writer, data)
	if err != nil {
		return "", err
	}
	return string(writer.Bytes()), nil
}

type GoTemplateEngine struct {
	FuncMap template.FuncMap
}

func (e GoTemplateEngine) Compile(viewFiles ...string) (render.CompiledView, error) {
	t := template.New(filepath.Base(viewFiles[0]))
	t.Funcs(e.FuncMap)
	_, err := t.ParseFiles(viewFiles...)
	if err != nil {
		return nil, err
	}
	tv := GoTemplateView(*t)
	return &tv, nil
}

var Engine = New()
