package requestextrator

import "net/http"

type RequestExtractor interface {
	ExtractRequest(r *http.Request) ([]byte, error)
}

type Formatter interface {
	Format([]byte) ([]byte, bool, error)
}

type ExtratorField struct {
	Extrator   RequestExtractor
	Formatters []Formatter
}

func (i *ExtratorField) Load(r *http.Request) ([]byte, bool, error) {
	var info []byte
	var ok bool
	var err error
	info, err = i.Extrator.ExtractRequest(r)
	if err != nil {
		return nil, false, err
	}
	for k := range i.Formatters {
		info, ok, err = i.Formatters[k].Format(info)
		if ok == false || err != nil {
			return nil, false, err
		}
	}
	return info, true, nil
}

func NewExtratorField() *ExtratorField {
	return &ExtratorField{}
}
