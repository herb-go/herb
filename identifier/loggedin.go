package identifier

import "net/http"

var LoggedInVerfier = IDVerifierFunc(func(id string) (bool, error) {
	return id != "", nil
})
var UnauthorizedHanlder = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(401), 401)
})

func NewLoggedInFilter(i Identifier, h http.Handler) *PlainFilter {
	g := NewFilter()
	g.Identifier = i
	g.IDVerifier = LoggedInVerfier
	if h == nil {
		g.Handler = UnauthorizedHanlder
	} else {
		g.Handler = h
	}
	return g
}
