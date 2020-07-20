package account

import (
	"github.com/herb-go/herbsecurity/authority"
)

type Account struct {
	authority.Authority
	authority.Passphrase
}

type Verified struct {
	*authority.Expiration
	authority.Principal
	authority.Agent
}
