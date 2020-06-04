package gotemplate

import (
	"html/template"
	"path"
	"path/filepath"

	"bytes"

	"github.com/herb-go/herb/ui/render"
)

func unescaped(x string) interface{} { return template.HTML(x) }

//New create new go template engine
func New() *RenderEngine {
	e := RenderEngine{
		FuncMap: template.FuncMap{
			"raw": unescaped,
		},
	}
	return &e
}

//View go template view
type View template.Template

//Execute execute view with given render data.
//Return render result as []byte and any error if raised.
func (v *View) Execute(data interface{}) ([]byte, error) {
	writer := bytes.NewBuffer([]byte{})
	t := template.Template(*v)
	err := t.Execute(writer, data)
	if err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}

//RenderEngine render engine main struct
type RenderEngine struct {
	//FuncMap builtin func map
	FuncMap template.FuncMap
	//ViewRoot view root path
	ViewRoot string
}

//SetViewRoot set view root path
func (e *RenderEngine) SetViewRoot(path string) {
	e.ViewRoot = path
}

//Compile complie view files to complied view.
func (e *RenderEngine) Compile(config *render.ViewConfig) (render.CompiledView, error) {
	var absFiles = make([]string, len(config.Files))

	for k, v := range config.Files {
		var p string
		if path.IsAbs(v) {
			p = v
		} else {
			p = path.Join(e.ViewRoot, v)
		}
		absFiles[k] = path.Clean(p)
	}
	t := template.New(filepath.Base(config.Files[0]))
	t.Funcs(e.FuncMap)
	_, err := t.ParseFiles(absFiles...)
	if err != nil {
		return nil, err
	}
	tv := View(*t)
	return &tv, nil
}

//RegisterFunc register func to engine
//Return any error if raised.
func (e *RenderEngine) RegisterFunc(name string, fn interface{}) error {
	e.FuncMap[name] = fn
	return nil
}

//Engine default go template render engine
var Engine = New()
