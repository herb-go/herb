package db

import (
	"encoding/json"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestDB(t *testing.T) {
	config := &DBConfig{}
	err := json.Unmarshal([]byte(ConfigJSON), config)
	if err != nil {
		t.Fatal(err)
	}
	db := New()
	err = db.Init(config)
	if err != nil {
		t.Fatal(err)
	}
	if db.Prefix() != config.Prefix {
		t.Error(db.Prefix())
	}
	dm := db.Table("test")
	if dm.DB() != db.DB() {
		t.Error("db not equal")
	}
	if dm.Name() != "test" {
		t.Error(dm.Name())
	}
	if dm.TableName() != db.Prefix()+dm.Name() {
		t.Error(dm.TableName())
	}
	dm.SetName("testnew")
	if dm.Name() != "testnew" {
		t.Error(dm.TableName())
	}
}
