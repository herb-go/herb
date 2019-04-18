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
	config.TokenLifetimeInHour = 1
	config.TokenMaxLifetimeInDay = 7
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

func TestCacheLifetimeInDay(t *testing.T) {
	var err error
	config := &StoreConfig{}
	config.DriverName = DriverNameCacheStore
	config.Cache.Driver = "syncmapcache"
	config.TokenLifetimeInDay = 1
	store := New()
	err = config.ApplyTo(store)
	if err != nil {
		panic(err)
	}
	if store.TokenLifetime != 24*time.Hour {
		t.Fatal(store.TokenLifetime)
	}

}
func TestClientStoreConfig(t *testing.T) {
	var err error
	config := &StoreConfig{}
	config.DriverName = DriverNameClientStore
	config.ClientStoreKey = "test"
	config.TokenLifetimeInHour = 1
	config.TokenMaxLifetimeInDay = 7
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
