package httpuser

import "net/http"

//CredentialProvider credential provider
type CredentialProvider interface {
	// credential request with given id
	CredentialRequest(r *http.Request) error
}
