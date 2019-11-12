package guarder

import (
	"net/http"

	"github.com/herb-go/herb/service"
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

func createBasicAuthWithConfig(conf service.Config, prefix string) (*BasicAuth, error) {
	var err error
	v := NewBasicAuth()
	return v, err
}

func basicAuthMapperFactory(conf service.Config, prefix string) (Mapper, error) {
	return createBasicAuthWithConfig(conf, prefix)
}

func registerBasicAuthFactory() {
	RegisterMapper("basicauth", basicAuthMapperFactory)
}

func init() {
	registerBasicAuthFactory()
}
