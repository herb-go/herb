package protecter

type Context struct {
	ID        string
	Protecter *Protecter
}

func NewContext() *Context {
	return &Context{}
}
