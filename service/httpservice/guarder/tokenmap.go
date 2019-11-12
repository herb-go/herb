package guarder

import (
	"strings"

	"github.com/herb-go/herb/service"
)

func NewTokenMap() *TokenMap {
	return &TokenMap{}
}

type TokenMap struct {
	ToLower bool
	Tokens  map[string]string
}

func (t *TokenMap) IdentifyParams(p *Params) (string, error) {
	id := p.ID()
	if id == "" {
		return "", nil
	}
	if t.ToLower {
		id = strings.ToLower(id)
	}
	token := t.Tokens[id]
	if token == "" || token != p.Credential() {
		return "", nil
	}
	return id, nil
}

func createTokenMapWithConfig(conf service.Config, prefix string) (*TokenMap, error) {
	var err error
	v := NewTokenMap()
	if err != nil {
		return nil, err
	}
	err = conf.Get("ToLower", &v.ToLower)
	if err != nil {
		return nil, err
	}
	err = conf.Get("Tokens", &v.Tokens)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func tokenMapIdentifierFactory(conf service.Config, prefix string) (Identifier, error) {
	return createTokenMapWithConfig(conf, prefix)
}
func registerTokenMapFactory() {
	RegisterIdentifier("tokenmap", tokenMapIdentifierFactory)
}

func init() {
	registerTokenMapFactory()
}
