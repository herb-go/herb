package middleware

import (
	"net/http"
	"testing"
)

type responsewriter struct {
	StatusCode int
}

func (w *responsewriter) Header() http.Header {
	h := http.Header{}
	h.Set("test", "test")
	return h
}
func (w *responsewriter) Write([]byte) (int, error) {
	return 11, nil
}
func (w *responsewriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

type responsewriterflusher struct {
	responsewriter
}

func TestResponsewriter(t *testing.T) {
	rw := &responsewriter{
		StatusCode: 200,
	}
	w := WrapResponseWriter(rw)
	if w.(*WrappedResponseWriter) == nil {
		t.Fatal(w)
	}
	if w.Header().Get("test") != "test" {
		t.Fatal(w)
	}
	if p, err := w.Write(nil); p != 11 || err != nil {
		t.Fatal(w)
	}
	if rw.StatusCode != 200 {
		t.Fatal(rw)
	}
	rw.WriteHeader(22)
	if rw.StatusCode != 22 {
		t.Fatal(rw)
	}
}
