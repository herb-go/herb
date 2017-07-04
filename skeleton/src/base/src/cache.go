package main

import "github.com/herb-go/herb/cache"

var Cache = cache.New()

func InitCache() {
	Must(Cache.OpenConfig(AppConfig.Cache))
}
