package render

import (
	"encoding/json"
	"io/ioutil"
)

type Option interface {
	ApplyTo(*Renderer) error
}

type OptionFunc func(*Renderer) error

func (i OptionFunc) ApplyTo(r *Renderer) error {
	return i(r)
}
func OptionCommon(e Engine, viewRoot string) OptionFunc {
	return func(r *Renderer) error {
		r.engine = e
		e.SetViewRoot(viewRoot)
		return nil
	}
}

type ViewsOption interface {
	ApplyTo(*Renderer) (map[string]*NamedView, error)
}
type ViewsOptionFunc func(*Renderer) (map[string]*NamedView, error)

func (o ViewsOptionFunc) ApplyTo(r *Renderer) (map[string]*NamedView, error) {
	return o(r)
}
func ViewsOptionCommon(data map[string]ViewConfig) ViewsOptionFunc {
	return func(r *Renderer) (map[string]*NamedView, error) {
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
}

type ViewsConf string

func (o ViewsConf) ApplyTo(r *Renderer) (map[string]*NamedView, error) {
	var data map[string]ViewConfig
	bytes, err := ioutil.ReadFile(string(o))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return ViewsOptionCommon(data).ApplyTo(r)
}
