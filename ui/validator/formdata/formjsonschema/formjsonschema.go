package formjsonschema

import (
	"github.com/herb-go/herb/ui"
	"github.com/herb-go/gojsonschema"
	"github.com/herb-go/herb/ui/validator/formdata"
)

type Schema struct {
	schema *gojsonschema.Schema
}

func (s *Schema) ValidateJSON(data []byte) ([]*Result, error) {
	results, err := s.schema.Validate(gojsonschema.NewBytesLoader(data))
	if err != nil {
		return nil, err
	}
	rerrs := results.Errors()
	r := make([]*Result, len(rerrs))
	for k, v := range rerrs {
		r[k] = &Result{
			result: v,
		}
	}
	return r, nil
}

type Result struct {
	result gojsonschema.ResultError
}

func (r *Result) Field() string {
	return r.result.EncodedPointer()
}

func (r *Result) Label() string {
	return r.result.SchemaTitle()
}

func (r *Result) Message() string {
	return r.result.DescriptionFormat()
}

func (r *Result) MessageData() map[string]string {
	data := map[string]string{}
	for k, v := range r.result.Details() {
		str, ok := v.(string)
		if ok {
			data[k] = str
		}
	}
	return data
}

func (s *Schema) Validate(f formdata.RequestValidator, data []byte) error {
	rs, err := s.ValidateJSON(data)
	if err != nil {
		return err
	}
	for k := range rs {
		f.AddError(ui.Replace())
	}
	return nil
}
