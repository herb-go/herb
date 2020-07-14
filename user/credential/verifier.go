package credential

type Verifier interface {
	Verify(*Collection) (string, error)
	Depenencies() map[Type]bool
}

type PlainVerifier struct {
	verifyFunc  func(*Collection) (string, error)
	depenencies map[Type]bool
}

func (v *PlainVerifier) Verify(c *Collection) (string, error) {
	return v.verifyFunc(c)
}

func (v *PlainVerifier) Depenencies() map[Type]bool {
	return v.depenencies
}

func VerifierFunc(verifyFunc func(*Collection) (string, error), depenencies ...Type) Verifier {
	v := &PlainVerifier{
		verifyFunc:  verifyFunc,
		depenencies: map[Type]bool{},
	}
	for k := range depenencies {
		v.depenencies[depenencies[k]] = true
	}
	return v
}

func Verify(v Verifier, l ...Credential) (string, error) {
	d := NewCollection()
	availableTypes := v.Depenencies()
	for k := range l {
		t := l[k].Type()
		if availableTypes[t] {
			val, err := l[k].Data()
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

var ForbiddenVerifier = FixedVerifier("")
