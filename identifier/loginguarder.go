package identifier

import "net/http"

var LoginVerfier = IDVerifierFunc(func(id string) (bool, error) {
	return id != "", nil
})
var UnauthorizedHanlder = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(401), 401)
})

func NewLoginGuarder(i Identifier, h http.Handler) *PlainGuarder {
	g := NewGuarder()
	g.Identifier = i
	g.IDVerifier = LoginVerfier
	if h == nil {
		g.Handler = UnauthorizedHanlder
	} else {
		g.Handler = h
	}
	return g
}
