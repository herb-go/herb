package querybuilder

import "testing"

func TestQuery(t *testing.T) {
	var q *PlainQuery
	if q.QueryCommand() != "" {
		t.Fatal(q)
	}
	if len(q.QueryArgs()) != 0 {
		t.Fatal(q)
	}
	b := New()
	q = b.New("test")
	q.And(b.New("test2"))
	if q.QueryCommand() != "( test AND test2 )" {
		t.Fatal(q)
	}
	if len(q.QueryArgs()) != 0 {
		t.Fatal(q)
	}
	q.Or(b.New("test3"))
	if q.QueryCommand() != "( ( test AND test2 ) OR test3 )" {
		t.Fatal(q)
	}
	if len(q.QueryArgs()) != 0 {
		t.Fatal(q)
	}
	q = b.New("")
	q.And(b.New("test2"))
	if q.QueryCommand() != "test2" {
		t.Fatal(q)
	}
	if len(q.QueryArgs()) != 0 {
		t.Fatal(q)
	}
	q = b.New("")
	q.Or(b.New("test2"))
	if q.QueryCommand() != "test2" {
		t.Fatal(q)
	}
	if len(q.QueryArgs()) != 0 {
		t.Fatal(q)
	}

}
