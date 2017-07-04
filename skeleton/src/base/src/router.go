package main

import (
	"net/http"

	router "github.com/herb-go/herb/middleware-httprouter"
	misc "github.com/herb-go/herb/middleware-misc"
	"github.com/herb-go/herb/util"
)

func NewRounter() http.Handler {
	var Router = router.New()
	Router.StripPrefix("/public").HandleFunc(misc.ServeFolder(http.Dir(util.Resource("/public/"))))
	return Router
}
