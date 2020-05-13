package guarder

import (
	"net/http"

	"github.com/herb-go/herb/user/credential"
)

func DefaultOnFail(g *Guarder) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(g.FailStatusCode), g.FailStatusCode)
	}
}

var DefaultVerifier = credential.FixedVerifier("")

type Guarder struct {
	Key            Key
	Credentials    []CredentialFactory
	Verifier       credential.Verifier
	OnFail         http.HandlerFunc
	FailStatusCode int
}

func (g *Guarder) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	loaders := make([]credential.Loader, len(g.Credentials))
	for k := range g.Credentials {
		loaders[k] = g.Credentials[k].CreateLoader(r)
	}
	id, err := credential.VerifyLoaders(g.Verifier, loaders...)
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
		Verifier:       DefaultVerifier,
		FailStatusCode: 403,
	}
	g.OnFail = DefaultOnFail(g)
	return g
}
