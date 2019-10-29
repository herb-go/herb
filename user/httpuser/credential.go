package httpuser

import "net/http"

//CredentialProvider credential provider
type CredentialProvider interface {
	// credential request with given id
	Credential(id string, r *http.Request) error
}
