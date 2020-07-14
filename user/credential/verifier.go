package credential

type Verifier interface {
	Verify(*Collection) (string, error)
	Depenencies() map[Type]bool
}

func Verify(v Verifier, c ...Credential) (string, error) {
	d := NewCollection()
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

func (v FixedVerifier) Verify(*Collection) (string, error) {
	return string(v), nil
}
func (v FixedVerifier) Depenencies() map[Type]bool {
	return map[Type]bool{}
}

type VerifierFactory interface {
	CreateVerifier(func(v interface{}) error) (Verifier, error)
}

var ForbiddenVerifier = FixedVerifier("")
