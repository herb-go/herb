package render

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
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
	engine    Engine
	ViewFiles map[string][]string
	Views     map[string]CompiledView
	lock      sync.RWMutex
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
	result, err := r.WriteJSON(w, data, status)
	if err != nil {
		panic(err)
	}
	return result
}

func (r *Renderer) HTMLFile(w http.ResponseWriter, path string, status int) (int, error) {
	return HTMLFile(w, path, status)
}
func (r *Renderer) MustHTMLFile(w http.ResponseWriter, path string, status int) int {
	result, err := r.HTMLFile(w, path, status)
	if err != nil {
		panic(err)
	}
	return result

}
func (r *Renderer) WriteHTML(w http.ResponseWriter, data []byte, status int) (int, error) {
	return WriteHTML(w, data, status)
}
func (r *Renderer) MustWriteHTML(w http.ResponseWriter, data []byte, status int) int {
	result, err := r.WriteHTML(w, data, status)
	if err != nil {
		panic(err)
	}
	return result

}
func (r *Renderer) JSON(w http.ResponseWriter, data interface{}, status int) (int, error) {
	return JSON(w, data, status)
}
func (r *Renderer) MustJSON(w http.ResponseWriter, data interface{}, status int) int {
	result, err := r.JSON(w, data, status)
	if err != nil {
		panic(err)
	}
	return result

}

func (r *Renderer) Error(w http.ResponseWriter, status int) (int, error) {
	return Error(w, status)
}
func (r *Renderer) MustError(w http.ResponseWriter, status int) int {
	result, err := Error(w, status)
	if err != nil {
		panic(err)
	}
	return result

}
func (r *Renderer) view(name string) (CompiledView, error) {
	var err error
	var view CompiledView
	if r.ViewFiles == nil {
		return nil, NewViewError(name, ErrorViewNotExist)
	}
	if r.Views != nil {
		view = r.Views[name]
	}
	if view == nil {
		vf := r.ViewFiles[name]
		if vf == nil {
			return nil, NewViewError(name, ErrorViewNotExist)
		}
		view, err = r.engine.Compile(vf...)
		if err != nil {
			return nil, err
		}
		r.setView(name, view)
	}
	return view, nil
}
func (r *Renderer) setView(name string, view CompiledView) {
	if r.Views == nil {
		r.Views = map[string]CompiledView{}
	}
	r.Views[name] = view
}
func (r *Renderer) setViewFiles(name string, viewFiles []string) {
	if r.ViewFiles == nil {
		r.ViewFiles = map[string][]string{}
	}
	r.ViewFiles[name] = viewFiles
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
	var loadedNamedViews = make(map[string]*NamedView, len(data))
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.Views == nil {
		r.Views = make(map[string]CompiledView, len(data))
	}
	for k, v := range data {
		delete(r.Views, k)
		r.setViewFiles(k, v.Views)
		loadedNamedViews[k] = &NamedView{
			Name:     k,
			Renderer: r,
		}
	}
	return loadedNamedViews, nil
}
func (r *Renderer) MustLoadViews(configpath string) map[string]*NamedView {
	vs, err := r.LoadViews(configpath)
	if err != nil {
		panic(err)
	}
	return vs
}
func (r *Renderer) GetView(ViewName string) *NamedView {
	return &NamedView{
		Name:     ViewName,
		Renderer: r,
	}
}

func (r *Renderer) Execute(viewname string, data interface{}) ([]byte, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	var cv, err = r.view(viewname)
	if err != nil {
		return nil, NewViewError(viewname, err)
	}
	if cv == nil {
		return nil, NewViewError(viewname, ErrorViewNotExist)
	}
	return cv.Execute(data)

}
func (r *Renderer) NewView(ViewName string, viewFiles ...string) *NamedView {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.Views, ViewName)
	r.setViewFiles(ViewName, viewFiles)
	v := &NamedView{
		Name:     ViewName,
		Renderer: r,
	}
	return v
}

type CompiledView interface {
	Execute(data interface{}) ([]byte, error)
}

type Engine interface {
	SetViewRoot(string)
	Compile(viewFiles ...string) (CompiledView, error)
}
