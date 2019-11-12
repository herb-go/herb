package guarder

import (
	"github.com/herb-go/herb/service"
)

type Token struct {
	Token string
	ID    StaticID
}

func NewToken() *Token {
	return &Token{}
}
func (t *Token) IdentifyParams(p *Params) (string, error) {
	if !t.ID.IsEmpty() {
		id := p.ID()
		if id == "" {
			return "", nil
		}
		if !t.ID.Equal(id) {
			return "", nil
		}
	}
	token := t.Token
	if token == "" || token != p.Credential() {
		return "", nil
	}
	return t.ID.ID(), nil
}

func (t *Token) CredentialParams() (*Params, error) {
	p := NewParams()
	if !t.ID.IsEmpty() {
		p.SetID(t.ID.ID())
	}
	p.SetCredential(t.Token)
	return p, nil
}

func createTokenWithConfig(conf service.Config, prefix string) (*Token, error) {
	var err error
	t := NewToken()
	err = conf.Get("Token", &t.Token)
	if err != nil {
		return nil, err
	}
	err = conf.Get("ID", &t.ID)
	if err != nil {
		return nil, err
	}
	return t, nil
}
func tokenCredentialFactory(conf service.Config, prefix string) (Credential, error) {
	return createTokenWithConfig(conf, prefix)
}
func tokenIdentifierFactory(conf service.Config, prefix string) (Identifier, error) {
	return createTokenWithConfig(conf, prefix)
}
func registerTokenFactory() {
	RegisterCredential("token", tokenCredentialFactory)
	RegisterIdentifier("token", tokenIdentifierFactory)
}

func init() {
	registerTokenFactory()
}
