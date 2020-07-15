package token

type Revoker interface {
	Revoke(id ID) error
}
