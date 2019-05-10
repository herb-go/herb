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
	selectquery.Select.AddRaw(15)
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
	if cmds != "SELECT prefix testfield1 , testfield2 , ?\nFROM tablename , table2name as table2alias\nLEFT JOIN table2 AS t2 ON field1=field2\nWHERE 1=1\nORDER BY testfield1 ASC  , testfield2 DESC \nLIMIT ? OFFSET ? \nother" {
		t.Fatal(cmds)
	}
	args := selectquery.QueryArgs()

	if len(args) != 3 || args[0] != 15 || args[1] != 5 || args[2] != 10 {
		t.Fatal(args)
	}
}

func TestUsing(t *testing.T) {
	testfield1 := ""
	testfield2 := ""
	fields := NewFields()
	fields.Set("testfield1", &testfield1)
	fields.Set("testfield2", &testfield2)
	builder := NewBuilder()
	selectquery := builder.NewSelect()
	selectquery.Select.AddFields(fields)
	selectquery.From.Add("tablename")
	selectquery.From.AddAlias("table2alias", "table2name")

	selectquery.Join.InnerJoin().Using("field1")
	cmds := selectquery.QueryCommand()
	if cmds != "SELECT testfield1 , testfield2\nFROM tablename , table2name as table2alias\nINNER JOIN  USING (field1)" {
		t.Fatal(cmds)
	}
	args := selectquery.QueryArgs()

	if len(args) != 0 {
		t.Fatal(args)
	}
}

func TestRightJoin(t *testing.T) {
	testfield1 := ""
	testfield2 := ""
	fields := NewFields()
	fields.Set("testfield1", &testfield1)
	fields.Set("testfield2", &testfield2)
	builder := NewBuilder()
	selectquery := builder.NewSelect()
	selectquery.Select.AddFields(fields)
	selectquery.From.Add("tablename")
	selectquery.From.AddAlias("table2alias", "table2name")

	selectquery.Join.RightJoin().On(builder.New("field1=field2"))
	cmds := selectquery.QueryCommand()
	if cmds != "SELECT testfield1 , testfield2\nFROM tablename , table2name as table2alias\nRIGHT JOIN  ON field1=field2" {
		t.Fatal(cmds)
	}
	args := selectquery.QueryArgs()

	if len(args) != 0 {
		t.Fatal(args)
	}
}

func TestResult(t *testing.T) {
	builder := NewBuilder()
	fields := NewFields()
	fields.Set("id", nil).Set("body", nil)
	selectquery := builder.NewSelect()
	selectquery.Select.AddFields(fields)
	r := selectquery.Result()
	if len(r.Fields) != 2 || len(r.args) != 2 {
		t.Fatal(r)
	}
	r.BindFields(fields)
	if len(r.Fields) != 2 || len(r.args) != 2 {
		t.Fatal(r)
	}
	r.BindFields(fields)
	if len(r.Fields) != 2 || len(r.args) != 2 {
		t.Fatal(r)
	}
	notusedfields := NewFields().Set("notuserd", 0)
	r.BindFields(notusedfields)
	if len(r.Fields) != 2 || len(r.args) != 2 {
		t.Fatal(r)
	}
}
