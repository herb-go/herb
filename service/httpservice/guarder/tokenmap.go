package guarder

import (
	"strings"
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

func createTokenMap(loader func(interface{}) error) (*TokenMap, error) {
	var err error
	v := NewTokenMap()
	err = loader(v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func tokenMapIdentifierFactory(loader func(interface{}) error) (Identifier, error) {
	return createTokenMap(loader)
}
func registerTokenMapFactory() {
	RegisterIdentifier("tokenmap", tokenMapIdentifierFactory)
}

func init() {
	registerTokenMapFactory()
}
