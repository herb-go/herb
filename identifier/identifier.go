package identifier

import "net/http"

//Identifier http request identifier
type Identifier interface {
	//IdentifyRequest identify http request
	//return identification and any error if rasied.
	IdentifyRequest(r *http.Request) (string, error)
}

//IDFunc identifier func type
type IDFunc func(r *http.Request) (string, error)

//IdentifyRequest identify http request
//return identification and any error if rasied.
func (f IDFunc) IdentifyRequest(r *http.Request) (string, error) {
	return f(r)
}

//FixedIdentifier
type FixedIdentifier string

//NopIdentifier fixed identifier ""
var NopIdentifier = FixedIdentifier("")

//IdentifyRequest identify http request
//return identification and any error if rasied.
func (i FixedIdentifier) IdentifyRequest(r *http.Request) (string, error) {
	return string(i), nil
}

type IDVerifier interface {
	VerifyID(id string) (bool, error)
}

type IDVerifierFunc func(id string) (bool, error)

func (f IDVerifierFunc) VerifyID(id string) (bool, error) {
	return f(id)
}
