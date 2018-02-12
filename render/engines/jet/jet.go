package jet

import (
	"bytes"
	"io"
	"os"
	"path"

	"github.com/CloudyKit/jet"
	"github.com/herb-go/herb/render"
)

type View struct {
	template *jet.Template
}

func (v *View) Execute(data interface{}) ([]byte, error) {
	var err error
	writer := bytes.NewBuffer([]byte{})
	err = v.template.Execute(writer, nil, data)

	if err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}

type JetEngine struct {
	Set      *jet.Set
	ViewRoot string
}

func (e *JetEngine) AddGlobal(Name string, fn interface{}) {
	(*e.Set).AddGlobal(Name, fn)
}
func (e *JetEngine) Compile(viewFiles ...string) (render.CompiledView, error) {
	if len(viewFiles) > 1 {
		return nil, render.ErrTooManyViewFiles
	}
	var t, err = e.Set.GetTemplate(viewFiles[0])
	if err != nil {
		return nil, err
	}
	tv := View{t}
	return &tv, nil
}
func (e *JetEngine) SetViewRoot(path string) {
	e.ViewRoot = path
}

var Engine = New()

func New() *JetEngine {
	var e = &JetEngine{}
	e.Set = jet.NewHTMLSetLoader(newLoader(e))
	return e
}

type loader struct {
	engine *JetEngine
}

func (l *loader) Open(name string) (io.ReadCloser, error) {
	return os.Open(name)
}

func (l *loader) Exists(name string) (string, bool) {
	fileName := path.Join(l.engine.ViewRoot, name)
	if _, err := os.Stat(fileName); err == nil {
		return fileName, true
	}
	return "", false

}

func newLoader(e *JetEngine) *loader {
	return &loader{
		engine: e,
	}
}
