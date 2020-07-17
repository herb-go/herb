package token

import "github.com/herb-go/herb/secret/auth"

type User string
type Token struct {
	*auth.Entity
	Owner auth.Owner
	User  User
	*auth.Expiration
}
