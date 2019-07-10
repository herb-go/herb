package columns

import (
	"strings"
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
	Register("test", func() Loader {
		return driver
	})
	notexist, err := Driver("notexist")
	if err == nil || !strings.Contains(err.Error(), "notexist") || notexist != nil {
		t.Fatal()
	}
	d, err := Driver("test")
	if err != nil || d != driver {
		t.Fatal(err)
	}
}
