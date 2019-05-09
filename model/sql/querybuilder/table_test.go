package querybuilder_test

import (
	"encoding/json"
	"testing"

	"github.com/herb-go/herb/model/sql/db"
	"github.com/herb-go/herb/model/sql/querybuilder"
)

func TestAlias(t *testing.T) {
	var DB = db.New()
	var config = db.NewConfig()
	var err error
	err = json.Unmarshal([]byte(MysqlConfigJSON), config)
	if err != nil {
		t.Fatal(err)
	}
	err = DB.Init(config)
	if err != nil {
		t.Fatal(err)
	}
	table1 := querybuilder.NewTable(DB.Table("testtable1"))
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
