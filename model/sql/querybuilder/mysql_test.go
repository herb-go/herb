package querybuilder_test

import (
	"encoding/json"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/herb-go/herb/model/sql/db"
	"github.com/herb-go/herb/model/sql/querybuilder"
)

func TestMysql(t *testing.T) {
	type Result struct {
		ID   string
		Body string
	}
	querybuilder.Debug = true
	var err error
	var DB = db.New()
	var config = db.NewConfig()
	err = json.Unmarshal([]byte(ConfigJSON), config)
	if err != nil {
		t.Fatal(err)
	}
	err = DB.Init(config)
	if err != nil {
		t.Fatal(err)
	}
	table1 := querybuilder.NewTable(DB.Table("testtable1"))
	if table1.Driver() != "mysql" {
		t.Fatal(table1)
	}
	truncatequery := table1.QueryBuilder().New("truncate table testtable1")
	truncatequery.MustExec(table1)

	_, err = DB.Exec("truncate table testtable2")
	if err != nil {
		t.Fatal(err)
	}

	builder := table1.QueryBuilder()
	fields := querybuilder.NewFields()
	var count int
	fields.Set(table1.QueryBuilder().CountField(), &count)
	countquery := table1.BuildCount()
	r := countquery.QueryRow(table1)
	err = countquery.Result().BindFields(fields).ScanFrom(r)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal(err)
	}
	insertquery := table1.NewInsert()
	fields = querybuilder.NewFields()
	fields.Set("id", "testid").Set("body", "testbody")
	insertquery.Insert.AddFields(fields)
	insertquery.Other = builder.New("ON DUPLICATE KEY UPDATE body= ?", "testbodydup")
	_, err = insertquery.Query().Exec(table1)
	if err != nil {
		t.Fatal(err)
	}

}
func init() {

}
