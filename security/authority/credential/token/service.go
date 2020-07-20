package token

import (
	"github.com/herb-go/herbsecurity/authority"
)

type Service interface {
	Issuer
	Revoker
	Refresher
	Loader
}
type Issuer interface {
	IssueToken(authority.Principal, authority.Passphrase, *authority.Expiration, Payload) (Token, error)
}

type Revoker interface {
	RevokeToken(authority.Passphrase) (Token, error)
}

type Refresher interface {
	RefreshToken(authority.Passphrase, *authority.Expiration) error
}

type Loader interface {
	LoadToken(authority.Passphrase) (*Token, error)
}
