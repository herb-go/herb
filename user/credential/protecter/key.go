package protecter

import (
	"context"
	"net/http"

	"github.com/herb-go/herb/user/credential"
)

type Key string

func (k Key) store(r *http.Request, ctx *Context) {
	reqctx := context.WithValue(r.Context(), k, ctx)
	req := r.WithContext(reqctx)
	*r = *req
}

func (k Key) load(r *http.Request) *Context {
	v := r.Context().Value(k)
	ctx, ok := v.(*Context)
	if !ok {
		ctx = NewContext()
		k.store(r, ctx)
	}
	return ctx

}
func (k Key) StoreID(r *http.Request, id string) {
	ctx := k.load(r)
	ctx.ID = id
	k.store(r, ctx)
}

func (k Key) LoadID(r *http.Request) string {
	ctx := k.load(r)
	return ctx.ID
}

func (k Key) IdentifyRequest(r *http.Request) (string, error) {
	return k.LoadID(r), nil
}

func (k Key) StoreProtecter(r *http.Request, p *Protecter) {
	if p == nil {
		return
	}
	ctx := k.load(r)
	ctx.Protecter = p
	k.store(r, ctx)
}

func (k Key) LoadProtecter(r *http.Request) *Protecter {
	ctx := k.load(r)
	if ctx.Protecter != nil {
		return ctx.Protecter
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
	k.ProtectMiddleware(k.LoadProtecter(r))(w, r, next)
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

func (k Key) ProtectMiddleware(p *Protecter) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if p == nil {
			DefaultOnFail(w, r)
			return
		}
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

func ProtectWith(p *Protecter, h http.Handler) http.Handler {
	return DefaultKey.ProtectWith(p, h)
}

func LoadID(r *http.Request) string {
	return DefaultKey.LoadID(r)
}

func ProtectMiddleware(p *Protecter) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return DefaultKey.ProtectMiddleware(p)
}
