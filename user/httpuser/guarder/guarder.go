package guarder

import "net/http"

type Guarder struct {
	Field      Field
	Middleware func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

func (g *Guarder) IdentifyRequest(r *http.Request) (string, error) {
	return g.Field.LoadID(r), nil
}

func (g *Guarder) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	g.Middleware(w, r, func(w http.ResponseWriter, r *http.Request) {
		if g.Field.LoadID(r) == "" {
			http.Error(w, http.StatusText(403), 403)
			return
		}
		next(w, r)
	})
}
