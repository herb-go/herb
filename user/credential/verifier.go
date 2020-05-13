package credential

type Verifier interface {
	VerifyCretentials(*Data) (string, error)
	Depenencies() map[Type]bool
}

func VerifyLoaders(v Verifier, loaders ...Loader) (string, error) {
	d := NewData()
	availableTypes := v.Depenencies()
	for k := range loaders {
		t := loaders[k].CredentialType()
		if availableTypes[t] {
			val, err := loaders[k].LoadCredential()
			if err != nil {
				return "", err
			}
			d.Append(t, val)
		}
	}
	return v.VerifyCretentials(d)

}

type FixedVerifier string

func (v FixedVerifier) VerifyCretentials(*Data) (string, error) {
	return string(v), nil
}
func (v FixedVerifier) Depenencies() map[Type]bool {
	return map[Type]bool{}
}
