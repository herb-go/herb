package guarder

import (
	"net/http"
)

type BasicAuth struct {
}

func NewBasicAuth() *BasicAuth {
	return &BasicAuth{}
}
func (h *BasicAuth) ReadParamsFromRequest(r *http.Request) (*Params, error) {
	p := NewParams()
	id, token, ok := r.BasicAuth()
	if !ok {
		return p, nil
	}
	p.SetID(id)
	p.SetCredential(token)
	return p, nil
}
func (h *BasicAuth) WriteParamsToRequest(r *http.Request, p *Params) error {
	r.SetBasicAuth(p.ID(), p.Credential())
	return nil
}

func createBasicAuth(loader func(interface{}) error) (*BasicAuth, error) {
	var err error
	v := NewBasicAuth()
	return v, err
}

func basicAuthMapperFactory(loader func(interface{}) error) (Mapper, error) {
	return createBasicAuth(loader)
}

func registerBasicAuthFactory() {
	RegisterMapper("basicauth", basicAuthMapperFactory)
}

func init() {
	registerBasicAuthFactory()
}
