package errorpage

import (
	"context"
	"net/http"
)

type contextName string

const contextNameDisable = contextName("ErrorPageDisabled")

//New create new error page middleware.
func New() *ErrorPage {
	return &ErrorPage{
		statusHandlers: map[int]func(w http.ResponseWriter, r *http.Request, status int){},
		errorHandler:   nil,
		ignoredStatus:  map[int]bool{},
	}
}

//ErrorPage  error page middleware main struct
type ErrorPage struct {
	statusHandlers map[int]func(w http.ResponseWriter, r *http.Request, status int)
	errorHandler   func(w http.ResponseWriter, r *http.Request, status int)
	ignoredStatus  map[int]bool
}

func (e *ErrorPage) disable(r *http.Request) {
	ctx := context.WithValue(r.Context(), contextNameDisable, true)
	*r = *r.WithContext(ctx)
}

//MiddlewareDisable middleware which disable previous installed error page middleware
func (e *ErrorPage) MiddlewareDisable(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	e.disable(r)
	next(w, r)
}

//ServeMiddleware serve as middleware
func (e *ErrorPage) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	newWriter := e.newResponseWriter(w, r)
	next(&newWriter, r)
	if newWriter.matched != nil {
		newWriter.matched(w, r, newWriter.status)
		e.disable(r)
	}
}
func (e *ErrorPage) getStatusHandler(status int) func(w http.ResponseWriter, r *http.Request, status int) {
	if e.ignoredStatus[status] {
		return nil
	}
	if statusHandlers, ok := e.statusHandlers[status]; ok {
		return statusHandlers

	}
	if status > 399 && e.errorHandler != nil {
		return e.errorHandler
	}
	return nil
}

//OnError configure  default error page when statuscode >399 and statuscode <600
func (e *ErrorPage) OnError(f func(w http.ResponseWriter, r *http.Request, status int)) *ErrorPage {
	e.errorHandler = f
	return e
}

//OnStatus configure error page by status code.
func (e *ErrorPage) OnStatus(status int, f func(w http.ResponseWriter, r *http.Request, status int)) *ErrorPage {
	e.statusHandlers[status] = f
	return e
}

//IgnoreStatus configure ignore given status.
func (e *ErrorPage) IgnoreStatus(status int) *ErrorPage {
	e.ignoredStatus[status] = true
	return e
}

func (e *ErrorPage) newResponseWriter(w http.ResponseWriter, r *http.Request) errorResponseWriter {
	return errorResponseWriter{
		ResponseWriter: w,
		req:            r,
		status:         0,
		ErrorPage:      *e,
		matched:        nil,
	}
}

type errorResponseWriter struct {
	http.ResponseWriter
	req       *http.Request
	status    int
	ErrorPage ErrorPage
	matched   func(w http.ResponseWriter, r *http.Request, status int)
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
	d := w.req.Context().Value(contextNameDisable)
	disabled, ok := d.(bool)
	if ok == false || disabled == false {
		w.matched = w.ErrorPage.getStatusHandler(status)

		if w.matched != nil {
			return
		}
	}
	w.ResponseWriter.WriteHeader(status)
}
