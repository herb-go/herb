package httpidentifier

import (
	"net/http"

	"github.com/herb-go/herb/user/identifier"
)

var DefaultOnFail = func(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(403), 403)
}

var DefaultIdentifier = identifier.FixedIdentifier("")

type Protecter struct {
	Credentialers []Credentialer
	Identifier    identifier.Identifier
	OnFail        http.HandlerFunc
}

func (g *Protecter) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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
	DefaultKey.StoreID(r, id)
	next(w, r)

}

func (g *Protecter) IdentifyRequest(r *http.Request) (string, error) {
	return DefaultKey.LoadID(r), nil
}

func NewProtecter() *Protecter {
	g := &Protecter{
		Identifier: DefaultIdentifier,
	}
	g.OnFail = DefaultOnFail
	return g
}
