package identifier

import "net/http"

type Guarder interface {
	Identifier
	IDVerifier
	http.Handler
}

func ServeGuarder(g Guarder, w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	id, err := g.IdentifyRequest(r)
	if err != nil {
		panic(err)
	}
	ok, err := g.VerifyID(id)
	if err != nil {
		panic(err)
	}
	if ok {
		g.ServeHTTP(w, r)
		return
	}
	next(w, r)
}

func NewGuarderMiddleware(g Guarder) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		ServeGuarder(g, w, r, next)
	}
}

type PlainGuarder struct {
	Identifier
	IDVerifier
	http.Handler
}

func (g *PlainGuarder) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ServeGuarder(g, w, r, next)
}
func NewGuarder() *PlainGuarder {
	return &PlainGuarder{}
}
