package guarder

type Builder interface {
	Build(*Guarder) error
}
