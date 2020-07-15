package token

type Storer interface {
	Store(*Token) error
}
