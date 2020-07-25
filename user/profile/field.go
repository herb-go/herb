package profile

type Field struct {
	Name  string
	Value string
}

func NewField() *Field {
	return &Field{}
}
