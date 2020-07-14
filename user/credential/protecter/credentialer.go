package protecter

import (
	"net/http"

	"github.com/herb-go/herb/user/credential"
)

type Credentialer interface {
	Credential(r *http.Request) credential.Credential
}

type CredentialerFactory interface {
	CreateCredentialer(func(v interface{}) error) (Credentialer, error)
}

type Credential struct {
	Request *http.Request
	Loader  *CredentialLoader
}

func (c *Credential) Type() credential.Type {
	return c.Loader.CredentialType
}
func (c *Credential) Load() (credential.Data, error) {
	return c.Loader.LoaderFunc(c.Request)
}

type CredentialLoader struct {
	CredentialType credential.Type
	LoaderFunc     func(*http.Request) (credential.Data, error)
}

func (c *CredentialLoader) Credential(r *http.Request) credential.Credential {
	return &Credential{
		Request: r,
		Loader:  c,
	}
}
