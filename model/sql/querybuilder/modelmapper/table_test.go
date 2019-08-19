package modelmapper_test

import (
	"encoding/json"
	"testing"

	"github.com/herb-go/herb/model/sql/db"
	"github.com/herb-go/herb/model/sql/querybuilder/modelmapper"
	_ "github.com/mattn/go-sqlite3"
)

var SqliteConfigJSON = `
{
	"Driver": "sqlite3",
	"DataSource": "_test/sqlite.db",
	"Prefix": "",
	"MaxIdleConns":10,
	"ConnMaxLifetimeInSecond":3600,
	"MaxOpenConns":10
}
`

func TestAlias(t *testing.T) {
	var DB = db.New()
	var config = db.NewConfig()
	var err error
	err = json.Unmarshal([]byte(SqliteConfigJSON), config)
	if err != nil {
		t.Fatal(err)
	}
	err = DB.Init(config)
	if err != nil {
		t.Fatal(err)
	}
	table1 := modelmapper.New(DB.Table("testtable1"))
	table1.SetAlias("")
	field := table1.FieldAlias("id")
	if field != "id" {
		t.Fatal(field)
	}
	table1.SetAlias("test")
	field = table1.FieldAlias("id")
	if field != "test.id" {
		t.Fatal(field)
	}
}
