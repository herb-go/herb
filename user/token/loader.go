package token

type Loader interface {
	Load(id ID) (*Token, error)
}
