package guarder

import (
	"github.com/herb-go/herb/user/identifier"
)

type Config struct {
	FailStatusCode int
	Credentialers  []*CredentialerConfig
	Identifier     *identifier.Config
}

type CredentialerConfig struct {
	Driver string
	Config func(v interface{}) error `config:", lazyload"`
}
