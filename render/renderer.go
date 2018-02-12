package render

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
)

//ViewConfig view config struct.
type ViewConfig struct {
	Views []string
}

//New create new renderer
func New() *Renderer {
	r := Renderer{}
	r.engine = nil
	r.Views = map[string]CompiledView{}
	r.ViewFiles = map[string][]string{}
	return &r
}

//Renderer renderer main struct
type Renderer struct {
	engine Engine
	//ViewFiles view file info map.
	ViewFiles map[string][]string
	//Views complied views map.
	Views map[string]CompiledView
	lock  sync.RWMutex
}

//Engine return engine of renderer.
func (r *Renderer) Engine() Engine {
	return r.engine
}

//Init set engine to renderer.
func (r *Renderer) Init(e Engine, viewRoot string) {
	r.engine = e
	e.SetViewRoot(viewRoot)
}

//WriteJSON write json data to response.
//Return bytes length wrote and any error if raised.
func (r *Renderer) WriteJSON(w http.ResponseWriter, data []byte, status int) (int, error) {
	return WriteJSON(w, data, status)
}

//MustWriteJSON write json data to response.
//Return bytes length wrote.
//Panic if any error raised.
func (r *Renderer) MustWriteJSON(w http.ResponseWriter, data []byte, status int) int {
	result, err := r.WriteJSON(w, data, status)
	if err != nil {
		panic(err)
	}
	return result
}

//HTMLFile write content of given file to response as html.
//Return bytes length wrote and any error if raised.
func (r *Renderer) HTMLFile(w http.ResponseWriter, path string, status int) (int, error) {
	return HTMLFile(w, path, status)
}

//MustHTMLFile write content of given file to response as html.
//Return bytes length wrote.
//Panic if any error raised.
func (r *Renderer) MustHTMLFile(w http.ResponseWriter, path string, status int) int {
	result, err := r.HTMLFile(w, path, status)
	if err != nil {
		panic(err)
	}
	return result

}

//WriteHTML write html data to response.
//Return bytes length wrote and any error if raised.
func (r *Renderer) WriteHTML(w http.ResponseWriter, data []byte, status int) (int, error) {
	return WriteHTML(w, data, status)
}

//MustWriteHTML write html data to response.
//Return bytes length wrote.
//Panic if any error raised.
func (r *Renderer) MustWriteHTML(w http.ResponseWriter, data []byte, status int) int {
	result, err := r.WriteHTML(w, data, status)
	if err != nil {
		panic(err)
	}
	return result

}

//JSON marshal data as json and write to response
//Return bytes length wrote and any error if raised.
func (r *Renderer) JSON(w http.ResponseWriter, data interface{}, status int) (int, error) {
	return JSON(w, data, status)
}

//MustJSON marshal data as json and write to response
//Return bytes length wrote.
//Panic if any error raised.
func (r *Renderer) MustJSON(w http.ResponseWriter, data interface{}, status int) int {
	result, err := r.JSON(w, data, status)
	if err != nil {
		panic(err)
	}
	return result

}

//Error write a http error to response
//Return bytes length wrote.
//Panic if any error raised.
func (r *Renderer) Error(w http.ResponseWriter, status int) (int, error) {
	return Error(w, status)
}

//MustError write a http error to response
//Return bytes length wrote.
//Panic if any error raised.
func (r *Renderer) MustError(w http.ResponseWriter, status int) int {
	result, err := r.Error(w, status)
	if err != nil {
		panic(err)
	}
	return result

}
func (r *Renderer) view(name string) (CompiledView, error) {
	var err error
	var view CompiledView
	view = r.Views[name]
	if view == nil {
		vf, ok := r.ViewFiles[name]
		if ok == false {
			return nil, NewViewError(name, ErrViewNotExist)
		}
		view, err = r.engine.Compile(vf...)
		if err != nil {
			return nil, NewViewError(name, err)
		}
		r.setView(name, view)
	}
	return view, nil
}
func (r *Renderer) setView(name string, view CompiledView) {
	r.Views[name] = view
}
func (r *Renderer) setViewFiles(name string, viewFiles []string) {
	r.ViewFiles[name] = viewFiles
}

//LoadViews load views form given file path in json format.
//Return loaded views and any error if raised.
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
	var loadedNamedViews = make(map[string]*NamedView, len(data))
	r.lock.Lock()
	defer r.lock.Unlock()
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

//MustLoadViews load views form given file path in json format.
//Return loaded views.
//Panic if any error raised.
func (r *Renderer) MustLoadViews(configpath string) map[string]*NamedView {
	vs, err := r.LoadViews(configpath)
	if err != nil {
		panic(err)
	}
	return vs
}

//GetView get view by name.
func (r *Renderer) GetView(ViewName string) *NamedView {
	return &NamedView{
		Name:     ViewName,
		Renderer: r,
	}
}

//Execute execute view by name with given render data.
//Return render result as []byte and any error if raised.
func (r *Renderer) Execute(viewname string, data interface{}) ([]byte, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	var cv, err = r.view(viewname)
	if err != nil {
		return nil, err
	}
	return cv.Execute(data)

}

//NewView create new view by name with given view files.
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

//CompiledView complied view interface.
type CompiledView interface {
	//Execute execute view with given render data.
	//Return render result as []byte and any error if raised.
	Execute(data interface{}) ([]byte, error)
}

//Engine render engine
type Engine interface {
	//SetViewRoot set view root path
	SetViewRoot(string)
	//Compile complie view files to complied view.
	Compile(viewFiles ...string) (CompiledView, error)
}
