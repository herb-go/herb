package gotemplate

import (
	"html/template"
	"path"
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

func (v *GoTemplateView) Execute(data interface{}) ([]byte, error) {

	writer := bytes.NewBuffer([]byte{})
	t := template.Template(*v)
	err := t.Execute(writer, data)
	if err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}

type GoTemplateEngine struct {
	FuncMap  template.FuncMap
	ViewRoot string
}

func (e *GoTemplateEngine) SetViewRoot(path string) {
	e.ViewRoot = path
}
func (e *GoTemplateEngine) Compile(viewFiles ...string) (render.CompiledView, error) {
	var absFiles = make([]string, len(viewFiles))

	for k, v := range viewFiles {
		var p string
		if path.IsAbs(v) {
			p = v
		} else {
			p = path.Join(e.ViewRoot, v)
		}
		absFiles[k] = path.Clean(p)
	}
	t := template.New(filepath.Base(viewFiles[0]))
	t.Funcs(e.FuncMap)
	_, err := t.ParseFiles(absFiles...)
	if err != nil {
		return nil, err
	}
	tv := GoTemplateView(*t)
	return &tv, nil
}

var Engine = New()
