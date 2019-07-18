package querybuilder

import "testing"

func TestDelete(t *testing.T) {
	builder := New()
	query := builder.NewDeleteQuery("testtable")
	query.Other = builder.New("other")
	query.Delete.Prefix = builder.New("prefix")
	query.Where.Condition = builder.Equal("testfield", "t1")
	q := query.Query()
	cmd := q.Command
	if cmd != "DELETE prefix FROM testtable\nWHERE testfield = ?\nother" {
		t.Fatal(cmd)
	}
	args := q.Args
	if len(args) != 1 || args[0] != "t1" {
		t.Fatal(args)
	}
}
