package guarder

import (
	"net/http"

	"github.com/herb-go/herb/service"
)

type IDTokenHeaders struct {
	IDHeader    string
	TokenHeader string
}

func NewIDTokenHeaders() *IDTokenHeaders {
	return &IDTokenHeaders{}
}
func (h *IDTokenHeaders) ReadParamsFromRequest(r *http.Request) (*Params, error) {
	p := NewParams()
	if h.IDHeader != "" {
		p.SetID(r.Header.Get(h.IDHeader))
	}
	p.SetCredential(r.Header.Get(h.TokenHeader))
	return p, nil
}
func (h *IDTokenHeaders) WriteParamsToRequest(r *http.Request, p *Params) error {
	if h.IDHeader != "" {
		r.Header.Set(h.IDHeader, p.ID())
	}
	r.Header.Set(h.TokenHeader, p.Credential())
	return nil
}

func createIDTokenHeadersWithConfig(conf service.Config, prefix string) (*IDTokenHeaders, error) {
	var err error
	v := NewIDTokenHeaders()
	if err != nil {
		return nil, err
	}
	err = conf.Get("IDHeader", &v.IDHeader)
	if err != nil {
		return nil, err
	}
	err = conf.Get("TokenHeader", &v.TokenHeader)
	if err != nil {
		return nil, err
	}
	if v.TokenHeader == "" {
		v.TokenHeader = "token"
	}
	return v, nil
}

func idTokenHeadersMapperFactory(conf service.Config, prefix string) (Mapper, error) {
	return createIDTokenHeadersWithConfig(conf, prefix)
}

func registerIDTokenHeadersFactory() {
	RegisterMapper("header", idTokenHeadersMapperFactory)
}

func init() {
	registerIDTokenHeadersFactory()
}
