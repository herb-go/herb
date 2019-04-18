package render

//Option renderer init option interface.
type Option interface {
	ApplyTo(*Renderer) error
}

// NewOptionCommon create new renderer option
func NewOptionCommon() *OptionCommon {
	return &OptionCommon{}
}

//OptionCommon Common renderer option
type OptionCommon struct {
	//Engine render engine
	Engine Engine
	//ViewRoot root path of view
	ViewRoot string
}

// ApplyTo apply option to renderer
func (o *OptionCommon) ApplyTo(r *Renderer) error {
	r.engine = o.Engine
	r.engine.SetViewRoot(o.ViewRoot)
	return nil
}

//ViewsOption renderer views init option.
type ViewsOption interface {
	ApplyTo(*Renderer) (map[string]*NamedView, error)
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
