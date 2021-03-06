package router

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/herb-go/herb/middleware"
)

//ContextName context name type
type ContextName string

//ContextNameRouterParams router params context name
const ContextNameRouterParams = ContextName("routerParams")

//Router router interface
type Router interface {
	//ServeHTTP serve as http handler
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	//StripPrefix strip request prefix and server as a middleware app
	StripPrefix(path string) *middleware.App
	//SetNotFoundHandler set not found handler
	SetNotFoundHandler(http.Handler)
}

//Param router param.
//Param should be parsed from url by router.
type Param struct {
	Name  string
	Value string
}

//Params params collection.
type Params []Param

//Get get param value by name.
func (p *Params) Get(name string) string {
	if p == nil || len(*p) == 0 {
		return ""
	}
	for k := range *p {
		if (*p)[k].Name == name {
			return (*p)[k].Value
		}
	}
	return ""
}

//Set get param value by name.
func (p *Params) Set(name string, value string) {
	params := Params{
		Param{
			Name:  name,
			Value: value,
		},
	}
	if p != nil {
		params = append(params, (*p)...)
	}
	*p = params
}

//GetParams get params from http request.
//If params collection does not exist,new param collection will be create and add to request.
func GetParams(r *http.Request) *Params {
	var params *Params
	p := r.Context().Value(ContextNameRouterParams)
	if p != nil {
		params = (p).(*Params)
	}
	if params == nil {
		params = &Params{}
		ctx := context.WithValue(r.Context(), ContextNameRouterParams, params)
		*r = *r.WithContext(ctx)
	}
	return params
}

func NewStripPrefixMiddleware(prefix string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
		next(w, r2)
	}
}
