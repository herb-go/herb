package httperror

import (
	"context"
	"net/http"
)

type contextNameType string

const ContextName = contextNameType("HTTPErrorDisabled")

func New() *Handler {
	return &Handler{
		statusHandlers: map[int]func(w http.ResponseWriter, r *http.Request, status int){},
		errorHandler:   nil,
		ignoredStatus:  map[int]bool{},
	}
}

type Handler struct {
	statusHandlers map[int]func(w http.ResponseWriter, r *http.Request, status int)
	errorHandler   func(w http.ResponseWriter, r *http.Request, status int)
	ignoredStatus  map[int]bool
}

func (h *Handler) disable(r *http.Request) {
	ctx := context.WithValue(r.Context(), ContextName, true)
	*r = *r.WithContext(ctx)
}
func (h *Handler) MiddlewareDisable(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h.disable(r)
	next(w, r)
}
func (h *Handler) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	newWriter := h.newResponseWriter(w, r)
	next(&newWriter, r)
	if newWriter.matched != nil {
		newWriter.matched(w, r, newWriter.status)
		h.disable(r)
	}
}
func (h *Handler) Handlers() []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){h.ServeMiddleware}

}
func (h *Handler) getStatusHandler(status int) func(w http.ResponseWriter, r *http.Request, status int) {
	if _, ok := h.ignoredStatus[status]; ok == true {
		return nil
	}
	if statusHandlers, ok := h.statusHandlers[status]; ok {
		return statusHandlers

	}
	if status > 399 && h.errorHandler != nil {
		return h.errorHandler
	}
	return nil
}
func (h *Handler) OnError(f func(w http.ResponseWriter, r *http.Request, status int)) *Handler {
	h.errorHandler = f
	return h
}
func (h *Handler) OnStatus(status int, f func(w http.ResponseWriter, r *http.Request, status int)) *Handler {
	h.statusHandlers[status] = f
	return h
}
func (h *Handler) IgnoreStatus(status int) *Handler {
	h.ignoredStatus[status] = true
	return h
}

func (h *Handler) newResponseWriter(w http.ResponseWriter, r *http.Request) errorResponseWriter {
	return errorResponseWriter{
		ResponseWriter: w,
		req:            r,
		status:         0,
		handler:        *h,
		matched:        nil,
	}
}

type errorResponseWriter struct {
	http.ResponseWriter
	req     *http.Request
	status  int
	handler Handler
	matched func(w http.ResponseWriter, r *http.Request, status int)
}

func (w *errorResponseWriter) Write(bytes []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	if w.matched != nil {
		return 0, nil
	}
	return w.ResponseWriter.Write(bytes)
}
func (w *errorResponseWriter) WriteHeader(status int) {
	w.status = status
	d := w.req.Context().Value(ContextName)
	disabled, ok := d.(bool)
	if ok == false || disabled == false {
		w.matched = w.handler.getStatusHandler(status)

		if w.matched != nil {
			return
		}
	}
	w.ResponseWriter.WriteHeader(status)
}
