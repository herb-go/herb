package token

import (
	"github.com/herb-go/herbsecurity/authority"
)

type Payload []byte

type Token struct {
	authority.Passphrase
	authority.Principal
	authority.Agent
	*authority.Expiration
	Payload []byte
}
