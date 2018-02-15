package render

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
