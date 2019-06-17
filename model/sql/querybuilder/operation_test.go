package querybuilder

import "testing"

func TestOperation(t *testing.T) {
	builder := New()
	q1 := builder.Equal("field1", "t1")
	q2 := builder.Equal("field2", "t2")
	q := builder.Concat(q1, q2)
	cmd := q.QueryCommand()
	if cmd != "field1 = ? field2 = ?" {
		t.Fatal(cmd)
	}
	args := q.QueryArgs()
	if len(args) != 2 || args[0] != "t1" || args[1] != "t2" {
		t.Fatal(args)
	}
	q = builder.Concat(q1, nil, q2)
	cmd = q.QueryCommand()
	if cmd != "field1 = ? field2 = ?" {
		t.Fatal(cmd)
	}
	args = q.QueryArgs()
	if len(args) != 2 || args[0] != "t1" || args[1] != "t2" {
		t.Fatal(args)
	}

	q = builder.Comma(q1, q2)
	cmd = q.QueryCommand()
	if cmd != "field1 = ? , field2 = ?" {
		t.Fatal(cmd)
	}
	args = q.QueryArgs()
	if len(args) != 2 || args[0] != "t1" || args[1] != "t2" {
		t.Fatal(args)
	}
	q = builder.Lines(q1, q2)
	cmd = q.QueryCommand()
	if cmd != "field1 = ?\nfield2 = ?" {
		t.Fatal(cmd)
	}
	args = q.QueryArgs()
	if len(args) != 2 || args[0] != "t1" || args[1] != "t2" {
		t.Fatal(args)
	}
	q = builder.And(q1, q2)
	cmd = q.QueryCommand()
	if cmd != "( field1 = ? AND field2 = ? )" {
		t.Fatal(cmd)
	}
	args = q.QueryArgs()
	if len(args) != 2 || args[0] != "t1" || args[1] != "t2" {
		t.Fatal(args)
	}
	q = builder.And(q1)
	cmd = q.QueryCommand()
	if cmd != "field1 = ?" {
		t.Fatal(cmd)
	}
	args = q.QueryArgs()
	if len(args) != 1 || args[0] != "t1" {
		t.Fatal(args)
	}

	q = builder.Or(q1, q2)
	cmd = q.QueryCommand()
	if cmd != "( field1 = ? OR field2 = ? )" {
		t.Fatal(cmd)
	}
	args = q.QueryArgs()
	if len(args) != 2 || args[0] != "t1" || args[1] != "t2" {
		t.Fatal(args)
	}
	q = builder.Or(q1)
	cmd = q.QueryCommand()
	if cmd != "field1 = ?" {
		t.Fatal(cmd)
	}
	args = q.QueryArgs()
	if len(args) != 1 || args[0] != "t1" {
		t.Fatal(args)
	}
}
