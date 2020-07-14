package protecter

import (
	"context"
	"net/http"

	"github.com/herb-go/herb/user/credential"
)

type Key string

func (k Key) StoreID(r *http.Request, id string) {
	reqctx := context.WithValue(r.Context(), k, id)
	req := r.WithContext(reqctx)
	*r = *req
}

func (k Key) LoadID(r *http.Request) string {
	v := r.Context().Value(k)
	id, _ := v.(string)
	return id
}

func (k Key) IdentifyRequest(r *http.Request) (string, error) {
	return k.LoadID(r), nil
}

func (k Key) StoreProtecter(r *http.Request, p *Protecter) {
	if p == nil {
		return
	}
	ctx := context.WithValue((*r).Context(), k, p)
	*r = *r.WithContext(ctx)
}

func (k Key) LoadProtecter(r *http.Request) *Protecter {
	v := r.Context().Value(k)
	if v != nil {
		return v.(*Protecter)
	}
	return DefaultProtecter
}
func (k Key) StoreProtecterMiddleware(p *Protecter) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		k.StoreProtecter(r, p)
		next(w, r)
	}
}
func (k Key) ServerMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	NewMiddleware(k, k.LoadProtecter(r))(w, r, next)
}
func (k Key) Protect(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		k.ServerMiddleware(w, r, h.ServeHTTP)
	})
}
func (k Key) ProtectWith(p *Protecter, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		k.StoreProtecterMiddleware(p)(w, r, k.Protect(h).ServeHTTP)
	})
}

func NewMiddleware(k Key, p *Protecter) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		credentials := make([]credential.Credential, len(p.Credentialers))
		for k := range p.Credentialers {
			credentials[k] = p.Credentialers[k].CredentialRequest(r)
		}
		id, err := credential.Verify(p.Verifier, credentials...)
		if err != nil {
			panic(err)
		}
		if id == "" {
			p.OnFail.ServeHTTP(w, r)
			return
		}
		k.StoreID(r, id)
		next(w, r)
	}
}

var DefaultKey = Key("")
