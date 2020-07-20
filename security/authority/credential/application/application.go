package application

import (
	"github.com/herb-go/herbsecurity/authority"
)

type Application struct {
	authority.Authority
	authority.Passphrase
}

type Verified struct {
	authority.Principal
	authority.Agent
	authority.Expiration
	*Application
}
