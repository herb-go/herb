package identifier

import "net/http"

type Filter interface {
	Identifier
	IDVerifier
	http.Handler
}

func ServeFilter(f Filter, w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	id, err := f.IdentifyRequest(r)
	if err != nil {
		panic(err)
	}
	ok, err := f.VerifyID(id)
	if err != nil {
		panic(err)
	}
	if !ok {
		f.ServeHTTP(w, r)
		return
	}
	next(w, r)
}

type PlainFilter struct {
	Identifier
	IDVerifier
	http.Handler
}

func (f *PlainFilter) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ServeFilter(f, w, r, next)
}

func (f *PlainFilter) MergeIdentifier(i Identifier) *PlainFilter {
	f.Identifier = i
	return f
}
func (f *PlainFilter) MergeIDVerifier(i IDVerifier) *PlainFilter {
	f.IDVerifier = i
	return f
}
func (f *PlainFilter) MergeHandler(h http.Handler) *PlainFilter {
	f.Handler = h
	return f
}
func NewFilter() *PlainFilter {
	return &PlainFilter{}
}
