package protecter

import (
	"context"
	"net/http"
)

type Key string

func (f Key) StoreID(r *http.Request, id string) {
	reqctx := context.WithValue(r.Context(), f, id)
	req := r.WithContext(reqctx)
	*r = *req
}

func (f Key) LoadID(r *http.Request) string {
	v := r.Context().Value(f)
	return v.(string)
}

var DefaultKey = Key("identifier")
