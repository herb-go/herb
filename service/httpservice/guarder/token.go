package guarder

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

func createToken(loader func(interface{}) error) (*Token, error) {
	var err error
	t := NewToken()
	err = loader(t)
	if err != nil {
		return nil, err
	}
	return t, nil
}
func tokenCredentialFactory(loader func(interface{}) error) (Credential, error) {
	return createToken(loader)
}
func tokenIdentifierFactory(loader func(interface{}) error) (Identifier, error) {
	return createToken(loader)
}
func registerTokenFactory() {
	RegisterCredential("token", tokenCredentialFactory)
	RegisterIdentifier("token", tokenIdentifierFactory)
}

func init() {
	registerTokenFactory()
}
