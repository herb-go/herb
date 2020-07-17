package token

import (
	"github.com/herb-go/herb/secret"
	"github.com/herb-go/herb/secret/auth"
)

type Service interface {
	Issuer
	Revoker
	Reseter
	Refresher
	Loader
}
type Issuer interface {
	Issue(auth.Owner, User, *auth.Expiration) (*Token, error)
}

type Refresher interface {
	Refresh(secret.ID, *auth.Expiration) error
}
type Loader interface {
	Load(secret.ID) (*Token, error)
}

type Revoker interface {
	Revoke(secret.ID) error
}

type Reseter interface {
	Reset(*auth.Entity) error
}
