package render

import (
	"net/http"
	"sync"
)

//New create new renderer
func New() *Renderer {
	r := Renderer{}
	r.engine = nil
	r.CompiledViews = map[string]CompiledView{}
	r.Views = map[string]*ViewConfig{}
	return &r
}

//NewViewConfig create new view config with given view files
func NewViewConfig(files ...string) *ViewConfig {
	return &ViewConfig{
		Files: files,
	}
}

//ViewConfig view config struct.
type ViewConfig struct {
	DevelopmentMode bool
	Files           []string
}

//Renderer renderer main struct
type Renderer struct {
	engine Engine
	//ViewFiles view file info map.
	Views map[string]*ViewConfig
	//Views complied views map.
	CompiledViews map[string]CompiledView
	//Developing Developing mode.If set to true,All views will not be cached.
	DevelopmentMode bool
	lock            sync.RWMutex
}

//Init init renderer with option.
func (r *Renderer) Init(option Option) error {
	return option.ApplyTo(r)
}

//Engine return engine of renderer.
func (r *Renderer) Engine() Engine {
	return r.engine
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
	view = r.CompiledViews[name]
	if view == nil {
		vconfig, ok := r.Views[name]
		if ok == false {
			return nil, NewViewError(name, ErrViewNotExist)
		}
		view, err = r.engine.Compile(vconfig)
		if err != nil {
			return nil, NewViewError(name, err)
		}
		config := r.Views[name]
		if !(config.DevelopmentMode == true || r.DevelopmentMode == true) {
			r.setCompiledView(name, view)
		}
	}
	return view, nil
}
func (r *Renderer) setCompiledView(name string, view CompiledView) {
	r.CompiledViews[name] = view
}
func (r *Renderer) setViewConfig(name string, config *ViewConfig) {
	r.Views[name] = config
	r.CompiledViews[name] = nil
}

//InitViews init renderer views with views option.
//Return inited views and any error if raised.
func (r *Renderer) InitViews(option ViewsOption) (map[string]*NamedView, error) {
	return option.ApplyTo(r)
}

//MustInitViews init renderer views with views option.
//Return inited views.
//Panic if any error raised.
func (r *Renderer) MustInitViews(option ViewsOption) map[string]*NamedView {
	vs, err := r.InitViews(option)
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
func (r *Renderer) NewView(ViewName string, config *ViewConfig) *NamedView {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.Views, ViewName)
	r.setViewConfig(ViewName, config)
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
	Compile(config *ViewConfig) (CompiledView, error)
	//RegisterFunc register func to engine
	//Return any error if raised.
	RegisterFunc(name string, fn interface{}) error
}
