package querybuilder

import (
	"testing"
)

func TestSelect(t *testing.T) {
	testfield1 := ""
	testfield2 := ""
	fields := NewFields()
	fields.Set("testfield1", &testfield1)
	fields.Set("testfield2", &testfield2)
	builder := NewBuilder()
	selectquery := builder.NewSelect()
	selectquery.Select.Prefix = builder.New("prefix")
	selectquery.Select.AddFields(fields)
	selectquery.From.Add("tablename")
	selectquery.From.AddAlias("table2alias", "table2name")
	selectquery.Limit.SetOffset(10)
	selectquery.Limit.SetLimit(5)
	selectquery.OrderBy.Add("testfield1", true)
	selectquery.OrderBy.Add("testfield2", false)
	selectquery.Join.LeftJoin().On(builder.New("field1=field2")).Alias("t2", "table2")
	selectquery.Where.Condition = builder.New("1=1")
	selectquery.Other = builder.New("other")
	cmds := selectquery.QueryCommand()
	if cmds != "SELECT prefix testfield1 , testfield2\nFROM tablename , table2name as table2alias\nLEFT Join table2 AS t2 ON field1=field2\nWHERE 1=1\nORDER BY testfield1 ASC  , testfield2 DESC \nLIMIT ? OFFSET ? \nother" {
		t.Fatal(cmds)
	}
	args := selectquery.QueryArgs()

	if len(args) != 2 || args[0] != 5 || args[1] != 10 {
		t.Fatal(args)
	}
}
