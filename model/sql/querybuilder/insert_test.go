package querybuilder

import "testing"

func TestInsert(t *testing.T) {
	fields := NewFields()
	fields.Set("field1", "t1")
	fields.Set("field2", "t2")
	builder := New()
	selectquery := builder.NewSelectQuery()
	selectquery.Select.AddFields(fields)
	selectquery.From.AddAlias("tb2", "table2")
	query := builder.NewInsertQuery("testtable")
	query.Insert.Prefix = builder.New("prefix")
	query.Insert.AddFields(fields)
	query.Insert.AddRaw("rawfield", "raw")
	query.Insert.AddSelect("t2.field3", selectquery)

	query.Other = builder.New("other")
	q := query.Query()
	cmd := q.Command
	if cmd != "INSERT prefix INTO testtable (field1 , field2 , rawfield , t2.field3 ) VALUES ( ? , ? , raw , ( SELECT field1 , field2\nFROM table2 AS tb2 ) )\nother" {
		t.Fatal(cmd)
	}
	args := q.Args
	if len(args) != 2 || args[0] != "t1" || args[1] != "t2" {
		t.Fatal(args)
	}
}

func TestInsertSubquery(t *testing.T) {
	fields := NewFields()
	fields.Set("field1", "t1")
	fields.Set("field2", "t2")
	builder := New()
	selectquery := builder.NewSelectQuery()
	selectquery.Select.AddFields(fields)
	selectquery.From.AddAlias("table2", "testtable2")
	query := builder.NewInsertQuery("testtable")
	query.Insert.Prefix = builder.New("prefix")
	query.Insert.AddFields(fields)
	query.Insert.WithSelect(selectquery)
	query.Other = builder.New("other")
	q := query.Query()
	cmd := q.Command
	if cmd != "INSERT prefix INTO testtable (field1 , field2 )\nSELECT field1 , field2\nFROM testtable2 AS table2\nother" {
		t.Fatal(cmd)
	}
	args := q.Args
	if len(args) != 0 {
		t.Fatal(args)
	}

}
