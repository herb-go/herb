package session

import (
	"testing"
	"time"
)

func TestCacheStoreConfig(t *testing.T) {
	var err error
	config := &StoreConfig{}
	config.DriverName = DriverNameCacheStore
	config.Cache.Driver = "syncmapcache"
	config.Cache.TTL = 3600
	config.TokenLifetime = "1h"
	config.TokenMaxLifetime = "168h"
	config.TokenContextName = "token"
	config.CookieName = "cookiename"
	config.CookiePath = "/"
	config.CookieSecure = true
	config.UpdateActiveIntervalInSecond = 100
	config.TokenLength = 32
	config.TokenPrefixMode = PrefixModeRaw
	store := New()
	err = config.ApplyTo(store)
	if err != nil {
		panic(err)
	}
	if store.TokenLifetime != 1*time.Hour {
		t.Fatal(store.TokenLifetime)
	}
	if store.TokenMaxLifetime != 7*24*time.Hour {
		t.Fatal(store.TokenMaxLifetime)
	}
	if store.CookieName != "cookiename" {
		t.Fatal(store.CookieName)
	}
	if store.CookiePath != "/" {
		t.Fatal(store.CookiePath)
	}
	if store.CookieSecure != true {
		t.Fatal(store.CookieSecure)
	}
	if store.UpdateActiveInterval != 100*time.Second {
		t.Fatal(store.UpdateActiveInterval)
	}
	if store.Driver.(*CacheDriver).Length != 32 {
		t.Fatal(store.Driver.(*CacheDriver).Length)
	}
	if store.Driver.(*CacheDriver).PrefixMode != PrefixModeRaw {
		t.Fatal(store.Driver.(*CacheDriver).Length)
	}
}

func TestClientStoreConfig(t *testing.T) {
	var err error
	config := &StoreConfig{}
	config.DriverName = DriverNameClientStore
	config.ClientStoreKey = "test"
	config.TokenLifetime = "1h"
	config.TokenMaxLifetime = "168h"
	config.TokenContextName = "token"
	config.CookieName = "cookiename"
	config.CookiePath = "/"
	config.CookieSecure = true
	config.UpdateActiveIntervalInSecond = 100
	store := New()
	err = config.ApplyTo(store)
	if err != nil {
		panic(err)
	}
	if store.TokenLifetime != 1*time.Hour {
		t.Fatal(store.TokenLifetime)
	}
	if store.TokenMaxLifetime != 7*24*time.Hour {
		t.Fatal(store.TokenMaxLifetime)
	}
	if store.CookieName != "cookiename" {
		t.Fatal(store.CookieName)
	}
	if store.CookiePath != "/" {
		t.Fatal(store.CookiePath)
	}
	if store.CookieSecure != true {
		t.Fatal(store.CookieSecure)
	}
	if store.UpdateActiveInterval != 100*time.Second {
		t.Fatal(store.UpdateActiveInterval)

	}
}

func TestCacheStoreDefaultConfig(t *testing.T) {
	var err error
	config := &StoreConfig{}
	config.DriverName = DriverNameCacheStore
	config.Cache.Driver = "syncmapcache"
	config.Cache.TTL = 3600
	store := New()
	err = config.ApplyTo(store)
	if err != nil {
		panic(err)
	}
	if store.TokenLifetime != defaultTokenLifetime {
		t.Fatal(store.TokenLifetime)
	}
	if store.TokenMaxLifetime != defaultTokenMaxLifetime {
		t.Fatal(store.TokenMaxLifetime)
	}
	if store.CookieName != defaultCookieName {
		t.Fatal(store.CookieName)
	}
	if store.CookiePath != defaultCookiePath {
		t.Fatal(store.CookiePath)
	}
	if store.CookieSecure != false {
		t.Fatal(store.CookieSecure)
	}
	if store.UpdateActiveInterval != defaultUpdateActiveInterval {
		t.Fatal(store.UpdateActiveInterval)

	}
}
