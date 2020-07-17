package auth

import "github.com/herb-go/herb/secret"

type Entity struct {
	secret.ID
	secret.Secret
}
