package querybuilder

import "testing"

func TestMisc(t *testing.T) {
	builder := NewBuilder()
	q := builder.NewValueList("t1", "t2")
	cmd := q.QueryCommand()
	if cmd != "? , ?" {
		t.Fatal(cmd)
	}
	args := q.QueryArgs()
	if len(args) != 2 || args[0] != "t1" || args[1] != "t2" {
		t.Fatal(args)
	}
	q = builder.NewValueList()
	cmd = q.QueryCommand()
	if cmd != "" {
		t.Fatal(cmd)
	}
	args = q.QueryArgs()
	if len(args) != 0 {
		t.Fatal(args)
	}
	q = builder.In("testfield", []string{"t1", "t2"})
	cmd = q.QueryCommand()
	if cmd != "testfield IN ( ? , ? )" {
		t.Fatal(cmd)
	}
	args = q.QueryArgs()
	if len(args) != 2 || args[0] != "t1" || args[1] != "t2" {
		t.Fatal(args)
	}
	q = builder.Equal("testfield", "t1")
	cmd = q.QueryCommand()
	if cmd != "testfield = ?" {
		t.Fatal(cmd)
	}
	args = q.QueryArgs()
	if len(args) != 1 || args[0] != "t1" {
		t.Fatal(args)
	}
	q = builder.Search("", "")
	cmd = q.QueryCommand()
	if cmd != "" {
		t.Fatal(cmd)
	}
	args = q.QueryArgs()
	if len(args) != 0 {
		t.Fatal(args)
	}
	q = builder.Search("testfield", "t1")
	cmd = q.QueryCommand()
	if cmd != "testfield LIKE ?" {
		t.Fatal(cmd)
	}
	args = q.QueryArgs()
	if len(args) != 1 || args[0] != "%t1%" {
		t.Fatal(args)
	}
	q = builder.Search("testfield", "\\test_%\\_\\%")
	cmd = q.QueryCommand()
	if cmd != "testfield LIKE ?" {
		t.Fatal(cmd)
	}
	args = q.QueryArgs()
	if len(args) != 1 || args[0] != "%\\\\test\\_\\%\\\\\\_\\\\\\%%" {
		t.Fatal(args)
	}

}