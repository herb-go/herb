package errorpage

import (
	"context"
	"net/http"

	"github.com/herb-go/herb/middleware"
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
	ctx := e.NewContext(w, r)
	next(ctx.NewWriter(), r)
	if ctx.matched != nil {
		ctx.matched(w, r, ctx.statusCode)
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
	if status >= 400 && e.errorHandler != nil {
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

//NewContext create new errorpage context with given writer and request
func (e *ErrorPage) NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		writer:    w,
		req:       r,
		ErrorPage: e,
	}
}

//Context errorpage context
type Context struct {
	writer     http.ResponseWriter
	req        *http.Request
	statusCode int
	ErrorPage  *ErrorPage
	matched    func(w http.ResponseWriter, r *http.Request, status int)
}

//Write Write response data
func (c *Context) Write(bytes []byte) (int, error) {
	if c.statusCode == 0 {
		c.writer.WriteHeader(http.StatusOK)
	}
	if c.matched != nil {
		return 0, nil
	}
	return c.writer.Write(bytes)
}

//WriteHeader writer header
func (c *Context) WriteHeader(statusCode int) {
	c.statusCode = statusCode
	d := c.req.Context().Value(contextNameDisable)
	disabled, ok := d.(bool)
	if ok == false || disabled == false {
		c.matched = c.ErrorPage.getStatusHandler(statusCode)
		if c.matched != nil {
			return
		}
	}
	c.writer.WriteHeader(statusCode)
}

//NewWriter create new response writer.
func (c *Context) NewWriter() http.ResponseWriter {
	w := middleware.WrapResponseWriter(c.writer)
	f := w.Functions()
	f.WriteFunc = c.Write
	f.WriteHeaderFunc = c.WriteHeader
	return w
}
