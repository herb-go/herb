package main

import (
	"net/http"

	"github.com/herb-go/herb/cache"
	"github.com/herb-go/herb/middleware"
	"github.com/herb-go/herb/resource"
	"github.com/herb-go/herb/util"
)

var Must = util.Must
var App = middleware.New()

type Config struct {
	NetConfig util.NetConfig
	Server    http.Server
	Cache     cache.Config
}

var AppConfig Config

func LoadConfigs() {
	resource.MustLoadJson(util.Config("config.json"), &AppConfig)
}
func InitModules() {
	InitCache()
	InitModel()

}
func main() {
	var Server = &AppConfig.Server
	LoadConfigs()
	InitModules()
	App.Use().Handle(NewRounter())
	defer util.Quit()
	util.MustListenAndServeHTTP(Server, AppConfig.NetConfig, App)
	util.WaitingQuit()
	util.ShutdownHTTP(Server)
}
