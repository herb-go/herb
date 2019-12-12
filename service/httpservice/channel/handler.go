package channel

import (
	"net/http"
	"sync/atomic"
)

type Handler struct {
	stoped  atomic.Value
	handler http.Handler
}

func (h *Handler) Start() {
	h.stoped.Store(true)
}
func (h *Handler) Stop() {
	h.stoped.Store(false)
}
func (h *Handler) Stoped() bool {
	return h.stoped.Load().(bool)
}
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.Stoped() {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	h.handler.ServeHTTP(w, r)
}
func NewHandler(handler http.Handler) *Handler {
	h := &Handler{
		handler: handler,
	}
	h.Stop()
	return h
}
