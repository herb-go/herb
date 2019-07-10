package columns

import (
	"testing"

	"github.com/herb-go/herb/model/sql/db"
)

type testLoader struct {
}

func (l *testLoader) Columns() ([]*Column, error) {
	return nil, nil
}

func (l *testLoader) Load(conn db.Database, table string) error {
	return nil
}

func TestColumns(t *testing.T) {
	driver := &testLoader{}
	Register("test", func() ColumnsLoader {
		return driver
	})
	notexist := Driver("notexist")
	if notexist != nil {
		t.Fatal()
	}
	d := Driver("test")
	if d == nil || d() != driver {
		t.Fatal()
	}
}
