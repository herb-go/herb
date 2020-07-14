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

func (p *Protecter) WithOnFail(h http.Handler) *Protecter {
	p.OnFail = h
	return p
}

func (p *Protecter) WithCredentialers(c ...Credentialer) *Protecter {
	p.Credentialers = c
	return p
}

func (p *Protecter) WithVerifier(v credential.Verifier) *Protecter {
	p.Verifier = v
	return p
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
