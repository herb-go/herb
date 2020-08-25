package identifier

import "net/http"

//Identifier http request identifier
type Identifier interface {
	//IdentifyRequest identify http request
	//return identification and any error if rasied.
	IdentifyRequest(r *http.Request) (string, error)
}
