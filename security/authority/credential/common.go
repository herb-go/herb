package credential

import (
	"github.com/herb-go/herbsecurity/authority"
)

var NamePassphrase = Name("passphrase")

func LoadPassphrase(c Credentials) (authority.Passphrase, error) {
	p, err := c.Get(NamePassphrase)
	if err != nil {
		return "", err
	}
	return authority.Passphrase(p), nil
}

var NameAuthority = Name("authority")

func LoadAuthority(c Credentials) (authority.Authority, error) {
	p, err := c.Get(NameAuthority)
	if err != nil {
		return "", err
	}
	return authority.Authority(p), nil
}
