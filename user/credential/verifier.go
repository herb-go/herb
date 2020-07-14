package credential

type Verifier interface {
	Verify(*CredentialDataCollection) (string, error)
	Depenencies() map[CredentialType]bool
}

func Verify(v Verifier, c ...Credential) (string, error) {
	d := NewDataCollection()
	availableTypes := v.Depenencies()
	for k := range c {
		t := c[k].Type()
		if availableTypes[t] {
			val, err := c[k].Load()
			if err != nil {
				return "", err
			}
			d.Set(t, val)
		}
	}
	return v.Verify(d)

}

type FixedVerifier string

func (v FixedVerifier) Verify(*CredentialDataCollection) (string, error) {
	return string(v), nil
}
func (v FixedVerifier) Depenencies() map[CredentialType]bool {
	return map[CredentialType]bool{}
}

type VerifierFactory interface {
	CreateVerifier(func(v interface{}) error) (Verifier, error)
}

var ForbiddenVerifier = FixedVerifier("")
