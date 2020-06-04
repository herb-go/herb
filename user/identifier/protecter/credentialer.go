package protecter

import (
	"net/http"

	"github.com/herb-go/herb/user/identifier"
)

type Credentialer interface {
	Credential(r *http.Request) identifier.Credential
}

type CredentialerFactory interface {
	CreateCredentialer(func(v interface{}) error) (Credentialer, error)
}

type Credential struct {
	Request *http.Request
	Loader  *CredentialLoader
}

func (c *Credential) Type() identifier.CredentialType {
	return c.Loader.CredentialType
}
func (c *Credential) Load() (identifier.CredentialData, error) {
	return c.Loader.LoaderFunc(c.Request)
}

type CredentialLoader struct {
	CredentialType identifier.CredentialType
	LoaderFunc     func(*http.Request) (identifier.CredentialData, error)
}

func (c *CredentialLoader) Credential(r *http.Request) identifier.Credential {
	return &Credential{
		Request: r,
		Loader:  c,
	}
}
