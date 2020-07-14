package protecter

import (
	"context"
	"net/http"
)

type Key string

func (k Key) StoreID(r *http.Request, id string) {
	reqctx := context.WithValue(r.Context(), k, id)
	req := r.WithContext(reqctx)
	*r = *req
}

func (k Key) LoadID(r *http.Request) string {
	v := r.Context().Value(k)
	return v.(string)
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
func (k Key) ProtecterMiddleware(p *Protecter) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		k.StoreProtecter(r, p)
		next(w, r)
	}
}
func (k Key) ServerMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	k.LoadProtecter(r).ServeMiddleware(w, r, next)
}
func (k Key) Protect(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		k.LoadProtecter(r).ServeMiddleware(w, r, h.ServeHTTP)
	})
}
func (k Key) ProtectWith(p *Protecter, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		k.ProtecterMiddleware(p)(w, r, k.Protect(h).ServeHTTP)
	})
}

var DefaultKey = Key("")
