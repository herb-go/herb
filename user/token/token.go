package token

type Token struct {
	ID     ID
	Secret Secret
}

func New() *Token {
	return &Token{}
}
