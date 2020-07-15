package token

type Map map[string][]byte

func (m *Map) Load(id ID) (*Token, error) {
	token := New()
	token.ID = id
	token.Secret = (*m)[string(id)]
	return token, nil
}

func (m *Map) Store(token *Token) error {
	(*m)[string(token.ID)] = token.Secret
	return nil
}

func (m *Map) Revoke(id ID) error {
	delete(*m, string(id))
	return nil
}

func NewMap() *Map {
	return &Map{}
}
