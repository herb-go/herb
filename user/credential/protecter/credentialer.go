package protecter

import (
	"net/http"

	"github.com/herb-go/herb/user/credential"
)

type Credentialer interface {
	CredentialRequest(r *http.Request) credential.Credential
}

type CredentialerFunc func(r *http.Request) credential.Credential

func (f CredentialerFunc) CredentialRequest(r *http.Request) credential.Credential {
	return f(r)
}
