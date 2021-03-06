package jet

import (
	"bytes"
	"io"
	"os"
	"path"

	"github.com/CloudyKit/jet"
	"github.com/herb-go/herb/ui/render"
)

//View jet template view
type View struct {
	template *jet.Template
}

//Execute execute view with given render data.
//Return render result as []byte and any error if raised.
func (v *View) Execute(data interface{}) ([]byte, error) {
	var err error
	writer := bytes.NewBuffer([]byte{})
	err = v.template.Execute(writer, nil, data)

	if err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}

//RenderEngine jet render engine main struct.
type RenderEngine struct {
	Set      *jet.Set
	ViewRoot string
}

//AddGlobal add builtin func.
func (e *RenderEngine) AddGlobal(Name string, fn interface{}) {
	(*e.Set).AddGlobal(Name, fn)
}

//Compile complie view files to complied view.
func (e *RenderEngine) Compile(config *render.ViewConfig) (render.CompiledView, error) {
	if len(config.Files) > 1 {
		return nil, render.ErrTooManyViewFiles
	}
	var t, err = e.Set.GetTemplate(config.Files[0])
	if err != nil {
		return nil, err
	}
	tv := View{t}
	return &tv, nil
}

//RegisterFunc register func to engine
//Return any error if raised.
func (e *RenderEngine) RegisterFunc(name string, fn interface{}) error {
	e.AddGlobal(name, fn)
	return nil
}

//SetViewRoot set view root path
func (e *RenderEngine) SetViewRoot(path string) {
	e.ViewRoot = path
}

//Engine default jet template render engine
var Engine = New()

//New create new jet template render engine.
func New() *RenderEngine {
	var e = &RenderEngine{}
	e.Set = jet.NewHTMLSetLoader(newLoader(e))
	e.Set.SetDevelopmentMode(true)
	return e
}

type loader struct {
	engine *RenderEngine
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

func newLoader(e *RenderEngine) *loader {
	return &loader{
		engine: e,
	}
}
