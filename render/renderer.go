package render

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type ViewConfig struct {
	Views []string
}

func New(e Engine, viewRoot string) *Renderer {
	r := Renderer{}
	r.engine = e
	if e != nil {
		r.engine.SetViewRoot(viewRoot)
	}
	r.Views = map[string]CompiledView{}
	return &r
}

type Renderer struct {
	engine Engine
	Views  map[string]CompiledView
}

func (r *Renderer) Engine() Engine {
	return r.engine
}
func (r *Renderer) SetEngine(e Engine, viewRoot string) {
	r.engine = e
	e.SetViewRoot(viewRoot)
}
func (r *Renderer) WriteJSON(w http.ResponseWriter, data []byte, status int) (int, error) {
	return WriteJSON(w, data, status)
}
func (r *Renderer) MustWriteJSON(w http.ResponseWriter, data []byte, status int) int {
	return MustWriteJSON(w, data, status)
}
func (r *Renderer) WriteHTML(w http.ResponseWriter, data []byte, status int) (int, error) {
	return r.WriteHTML(w, data, status)
}
func (r *Renderer) MustWriteHTML(w http.ResponseWriter, data []byte, status int) int {
	return r.MustWriteHTML(w, data, status)
}
func (r *Renderer) JSON(w http.ResponseWriter, data interface{}, status int) (int, error) {
	return JSON(w, data, status)
}
func (r *Renderer) MustJSON(w http.ResponseWriter, data interface{}, status int) int {
	return MustJSON(w, data, status)

}

func (r *Renderer) Error(w http.ResponseWriter, status int) (int, error) {
	return Error(w, status)
}
func (r *Renderer) MustError(w http.ResponseWriter, status int) int {
	return MustError(w, status)
}
func (r *Renderer) view(name string) CompiledView {
	if r.Views == nil {
		return nil
	}
	view, ok := r.Views[name]
	if ok == false {
		return nil
	}
	return view
}
func (r *Renderer) setView(name string, view CompiledView) {
	if r.Views == nil {
		r.Views = map[string]CompiledView{}
	}
	r.Views[name] = view
}
func (r *Renderer) LoadViews(configpath string) (map[string]*NamedView, error) {
	var data map[string]ViewConfig
	bytes, err := ioutil.ReadFile(configpath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	loadedViewsName := map[string]*NamedView{}
	for k, v := range data {
		namedView, err := r.NewView(k, v.Views...)
		if err != nil {
			return loadedViewsName, err
		}
		loadedViewsName[k] = namedView
	}
	return loadedViewsName, nil
}

func (r *Renderer) GetView(ViewName string) *NamedView {
	return &NamedView{
		Name:     ViewName,
		Renderer: r,
	}
}
func (r *Renderer) NewView(ViewName string, viewFiles ...string) (*NamedView, error) {
	cv, err := r.engine.Compile(viewFiles...)
	if err != nil {
		return nil, NewViewError(ViewName, err)
	}
	r.setView(ViewName, cv)
	v := &NamedView{
		Name:     ViewName,
		Renderer: r,
	}
	return v, nil
}

type CompiledView interface {
	Execute(data interface{}) (string, error)
}

type Engine interface {
	SetViewRoot(string)
	Compile(viewFiles ...string) (CompiledView, error)
}
