package protecter

import (
	"net/http"

	"github.com/herb-go/herb/user/credential"
)

var DefaultOnFail = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(403), 403)
})

var DefaultVerifier = credential.ForbiddenVerifier

type Protecter struct {
	Credentialers []Credentialer
	Verifier      credential.Verifier
	OnFail        http.Handler
}

func (p *Protecter) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	credentials := make([]credential.Credential, len(p.Credentialers))
	for k := range p.Credentialers {
		credentials[k] = p.Credentialers[k].CredentialRequest(r)
	}
	id, err := credential.Verify(p.Verifier, credentials...)
	if err != nil {
		panic(err)
	}
	if id == "" {
		p.OnFail.ServeHTTP(w, r)
		return
	}
	DefaultKey.StoreID(r, id)
	next(w, r)

}

func (g *Protecter) IdentifyRequest(r *http.Request) (string, error) {
	return DefaultKey.LoadID(r), nil
}

func New() *Protecter {
	p := &Protecter{
		Verifier: DefaultVerifier,
		OnFail:   DefaultOnFail,
	}
	return p
}

var ForbiddenProtecter = New()

var DefaultProtecter = ForbiddenProtecter

var NotWorkingProtecter = &Protecter{
	Verifier: credential.FixedVerifier("notworking"),
	OnFail:   DefaultOnFail,
}
