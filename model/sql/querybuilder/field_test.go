package querybuilder

import "testing"

func TestField(t *testing.T) {
	fields := NewFields()
	fields.Set("testfield1", "")
	fields.Set("testfield1", "t1")
	fields.Set("testfield2", "t2")
	if len(*fields) != 2 {
		t.Fatal(fields)
	}
}
