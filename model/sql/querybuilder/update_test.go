package querybuilder

import "testing"

func TestUpdate(t *testing.T) {
	fields := NewFields()
	fields.Set("field1", "t1")
	fields.Set("field2", "t2")
	builder := NewBuilder()
	selectquery := builder.NewSelect()
	selectquery.Select.AddFields(fields)
	selectquery.From.AddAlias("tb2", "table2")
	query := builder.NewUpdate("table")
	query.Update.AddFields(fields)
	query.Update.Prefix = builder.New("prefix")
	query.Update.AddRaw("raw", "raw")
	query.Update.SetAlias("testalias")
	query.Update.AddSelect("field3", selectquery)
	query.Where.Condition = builder.New("1=1")
	query.Other = builder.New("other")
	cmd := query.QueryCommand()
	if cmd != "UPDATE prefix table AS testalias SET field1 = ? , field2 = ? , raw = raw , field3 = ( SELECT field1 , field2\nFROM table2 as tb2 )\nWHERE 1=1\nother" {
		t.Fatal(cmd)
	}
	args := query.QueryArgs()
	if len(args) != 2 || args[0] != "t1" || args[1] != "t2" {
		t.Fatal(args)
	}
}