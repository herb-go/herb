package httpidentifier

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
