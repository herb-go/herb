package guarder

import "net/http"

type Guarder struct {
	Mapper     Mapper
	Identifier Identifier
}

func (g *Guarder) Init(o GuarderOption) error {
	return o.ApplyToGuarder(g)
}
func (g *Guarder) IdentifyRequest(r *http.Request) (string, error) {
	p, err := g.Mapper.ReadParamsFromRequest(r)
	if err != nil {
		return "", err
	}
	return g.Identifier.IdentifyParams(p)
}

func NewGuarder() *Guarder {
	return &Guarder{}
}
