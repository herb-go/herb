package extractor

import "net/http"
import "github.com/herb-go/herb/middleware/httpinfo"

type Extractor interface {
	Extract(r *http.Request) ([]byte, error)
}

type ExtractorFunc func(r *http.Request) ([]byte, error)

func (f ExtractorFunc) Extract(r *http.Request) ([]byte, error) {
	return f(r)
}

type ExtratorField struct {
	Extrator   Extractor
	Formatters []httpinfo.Formatter
}

func (i *ExtratorField) Load(r *http.Request) ([]byte, bool, error) {
	var info []byte
	var ok bool
	var err error
	info, err = i.Extrator.Extract(r)
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
