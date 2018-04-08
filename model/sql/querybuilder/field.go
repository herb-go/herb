package querybuilder

type Field struct {
	Field string
	Data  interface{}
}
type Fields []Field

func NewFields() *Fields {
	return &Fields{}
}
func (f Fields) Set(field string, data interface{}) Fields {
	for k, v := range f {
		if v.Field == field {
			f[k].Data = data
			return f
		}
	}
	f = append(f, Field{Field: field, Data: data})
	return f
}
