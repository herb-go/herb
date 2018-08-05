package render

//Option renderer init option interface.
type Option interface {
	ApplyTo(*Renderer) error
}

//OptionFunc renderer func init option.
type OptionFunc func(*Renderer) error

//ApplyTo apply option to renderer.
func (i OptionFunc) ApplyTo(r *Renderer) error {
	return i(r)
}

//OptionCommon option with given render engine and view root path.
func OptionCommon(e Engine, viewRoot string) OptionFunc {
	return func(r *Renderer) error {
		r.engine = e
		e.SetViewRoot(viewRoot)
		return nil
	}
}

//ViewsOption renderer views init option.
type ViewsOption interface {
	ApplyTo(*Renderer) (map[string]*NamedView, error)
}

//ViewsOptionFunc renderer views func init option.
type ViewsOptionFunc func(*Renderer) (map[string]*NamedView, error)

//ApplyTo apply views option to renderer.
func (o ViewsOptionFunc) ApplyTo(r *Renderer) (map[string]*NamedView, error) {
	return o(r)
}

//ViewsOptionCommon views option with new view configs.
type ViewsOptionCommon struct {
	DevelopmentMode bool
	Views           map[string]ViewConfig
}

//ApplyTo init renderer with given json conf.
func (o ViewsOptionCommon) ApplyTo(r *Renderer) (map[string]*NamedView, error) {
	var loadedNamedViews = make(map[string]*NamedView, len(o.Views))
	r.lock.Lock()
	defer r.lock.Unlock()
	for k, v := range o.Views {
		delete(r.Views, k)
		r.setViewConfig(k, v)
		loadedNamedViews[k] = &NamedView{
			Name:     k,
			Renderer: r,
		}
	}
	r.DevelopmentMode = o.DevelopmentMode
	return loadedNamedViews, nil
}
