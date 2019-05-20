package querybuilder

// Field query field struct
type Field struct {
	// Field field name
	Field string
	// Data field data
	Data interface{}
}

// Fields field list
type Fields []Field

// NewFields create new fields
func NewFields() *Fields {
	return &Fields{}
}

// Set set field value wth given field name dana.
// return fields self.
func (f *Fields) Set(field string, data interface{}) *Fields {
	for k, v := range *f {
		if v.Field == field {
			(*f)[k].Data = data
			return f
		}
	}
	*f = append(*f, Field{Field: field, Data: data})
	return f
}
