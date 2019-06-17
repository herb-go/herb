package querybuilder

// Field query field struct
type Field struct {
	// Field field name
	Field string
	// Data field data
	Data interface{}
}

//NewField Create new field
func NewField() *Field {
	return &Field{}
}

// Fields field list
type Fields []*Field

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
	datafield := NewField()
	datafield.Field = field
	datafield.Data = data
	*f = append(*f, datafield)
	return f
}
