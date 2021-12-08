package httpinfo

import "net/http"

type Extractor interface {
	Extract(r *http.Request) ([]byte, error)
}

type ExtractorFunc func(r *http.Request) ([]byte, error)

func (f ExtractorFunc) Extract(r *http.Request) ([]byte, error) {
	return f(r)
}

type ExtractorField struct {
	Extractor  Extractor
	Formatters []Formatter
}

func (f *ExtractorField) WithExtrator(e Extractor) *ExtractorField {
	f.Extractor = e
	return f
}
func (f *ExtractorField) WithFormatters(formatters ...Formatter) *ExtractorField {
	f.Formatters = formatters
	return f
}
func (f *ExtractorField) LoadInfo(r *http.Request) ([]byte, bool, error) {
	var info []byte
	var ok bool
	var err error
	info, err = f.Extractor.Extract(r)
	if err != nil {
		return nil, false, err
	}
	for k := range f.Formatters {
		info, ok, err = f.Formatters[k].Format(info)
		if ok == false || err != nil {
			return nil, false, err
		}
	}
	return info, true, nil
}
func (f *ExtractorField) IdentifyRequest(r *http.Request) (string, error) {
	data, ok, err := f.LoadInfo(r)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", nil
	}
	return string(data), nil
}

func NewExtractorField() *ExtractorField {
	return &ExtractorField{}
}
