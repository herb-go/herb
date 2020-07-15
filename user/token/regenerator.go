package token

type Regenerator struct {
	Generator Generator
	Storer    Storer
}

func (r *Regenerator) Regenerate(id ID) (*Token, error) {
	secret, err := r.Generator.Generate()
	if err != nil {
		return nil, err
	}
	t := New()
	t.ID = id
	t.Secret = secret
	err = r.Storer.Store(t)
	if err != nil {
		return nil, err
	}
	return t, nil
}
