package db

import (
	"encoding/json"
	"testing"
)

func TestConfig(t *testing.T) {
	config := NewConfig()
	err := json.Unmarshal([]byte(ConfigJSON), config)
	config.ConnMaxLifetimeInSecond = 10
	config.Driver = "mysql"
	config.Type = ""
	config.MaxIdleConns = 0
	config.MaxOpenConns = 0
	config.Prefix = "prefix"
	db := New()
	err = config.ApplyTo(db)
	if err != nil {
		t.Fatal(err)
	}
	if db.Driver() != config.Driver {
		t.Fatal(db.Driver())
	}
	config = NewConfig()
	err = json.Unmarshal([]byte(ConfigJSON), config)
	config.ConnMaxLifetimeInSecond = 0
	config.Driver = "mysql"
	config.Type = "othertype"
	config.MaxIdleConns = 0
	config.MaxOpenConns = 0
	config.Prefix = "prefix"
	db = New()
	err = config.ApplyTo(db)
	if err != nil {
		t.Fatal(err)
	}
	if db.Driver() != config.Type {
		t.Fatal(db.Driver())
	}
}
