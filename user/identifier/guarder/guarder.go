package guarder

import (
	"net/http"

	"github.com/herb-go/herb/user/identifier"
)

func DefaultOnFail(g *Guarder) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(g.FailStatusCode), g.FailStatusCode)
	}
}

var DefaultIdentifier = identifier.FixedIdentifier("")

type Guarder struct {
	Key            Key
	Credentialers  []Credentialer
	Identifier     identifier.Identifier
	OnFail         http.HandlerFunc
	FailStatusCode int
}

func (g *Guarder) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	credentials := make([]identifier.Credential, len(g.Credentialers))
	for k := range g.Credentialers {
		credentials[k] = g.Credentialers[k].Credential(r)
	}
	id, err := identifier.Identify(g.Identifier, credentials...)
	if err != nil {
		panic(err)
	}
	if id == "" {
		g.OnFail(w, r)
		return
	}
	g.Key.StoreID(r, id)
	next(w, r)

}

func (g *Guarder) IdentifyRequest(r *http.Request) (string, error) {
	return g.Key.LoadID(r), nil
}

func New() *Guarder {
	g := &Guarder{
		Identifier:     DefaultIdentifier,
		FailStatusCode: 403,
	}
	g.OnFail = DefaultOnFail(g)
	return g
}
