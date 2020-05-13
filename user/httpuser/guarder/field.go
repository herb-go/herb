package guarder

import (
	"context"
	"net/http"
)

type Field string

var DefaultField = Field("")

func (f Field) StoreID(r *http.Request, id string) {
	reqctx := context.WithValue(r.Context(), f, id)
	req := r.WithContext(reqctx)
	*r = *req
}

func (f Field) LoadID(r *http.Request) string {
	v := r.Context().Value(f)
	return v.(string)
}
