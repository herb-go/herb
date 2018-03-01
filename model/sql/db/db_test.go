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
	if dm.TableName() != "test" {
		t.Error(dm.TableName())
	}
	if dm.DBTableName() != db.Prefix()+dm.TableName() {
		t.Error(dm.DBTableName())
	}
	dm.SetTableName("testnew")
	if dm.TableName() != "testnew" {
		t.Error(dm.TableName())
	}
}
