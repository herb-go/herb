package account

import (
	"github.com/herb-go/herbsecurity/authority"
)

type Verifier interface {
	VerifyAccount(*Account) (Verified, error)
}

type Updater interface {
	UpdateAccount(authority.Principal, authority.Authority, authority.Passphrase) error
}

type Creator interface {
	CreateAccount(authority.Principal, authority.Authority, authority.Agent, authority.Passphrase) error
}
type Revoker interface {
	RevokeAccount(authority.Principal, authority.Authority) error
}
type Service interface {
	Verifier
	Updater
	Creator
}
