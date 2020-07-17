package auth

type Owner string

func (o Owner) OwnedBy() (Owner, error) {
	return o, nil
}
