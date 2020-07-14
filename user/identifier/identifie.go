package identifier

type Identifier interface {
	IdentifyCretentials(*CredentialDataCollection) (string, error)
	Depenencies() map[CredentialType]bool
}

func Identify(i Identifier, c ...Credential) (string, error) {
	d := NewDataCollection()
	availableTypes := i.Depenencies()
	for k := range c {
		t := c[k].Type()
		if availableTypes[t] {
			val, err := c[k].Load()
			if err != nil {
				return "", err
			}
			d.Append(t, val)
		}
	}
	return i.IdentifyCretentials(d)

}

type FixedIdentifier string

func (v FixedIdentifier) IdentifyCretentials(*CredentialDataCollection) (string, error) {
	return string(v), nil
}
func (v FixedIdentifier) Depenencies() map[CredentialType]bool {
	return map[CredentialType]bool{}
}

type IdentifierFactory interface {
	CreateIdentifier(func(v interface{}) error) (Identifier, error)
}

var ForbiddenIdentifier = FixedIdentifier("")
