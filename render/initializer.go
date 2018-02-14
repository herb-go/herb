package render

type Initializer interface {
	Init(*Renderer) error
}

type InitializerFunc func(*Renderer) error

func (i InitializerFunc) Init(r *Renderer) error {
	return i(r)
}
func Option(e Engine, viewRoot string) InitializerFunc {
	return func(r *Renderer) error {
		r.engine = e
		e.SetViewRoot(viewRoot)
		return nil
	}
}
