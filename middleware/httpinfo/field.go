package httpinfo

import (
	"net/http"
)

//Field request info field.
//Field is used to load specific information from http request
type Field interface {
	LoadInfo(r *http.Request) ([]byte, bool, error)
}

type FieldFunc func(r *http.Request) ([]byte, bool, error)

func (f FieldFunc) LoadInfo(r *http.Request) ([]byte, bool, error) {
	return f(r)
}

//StringField string field
type StringField struct {
	Field Field
}

func (f *StringField) LoadStringInfo(r *http.Request) (string, error) {
	data, _, err := f.Field.LoadInfo(r)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func NewStringField(f Field) *StringField {
	return &StringField{
		Field: f,
	}
}
