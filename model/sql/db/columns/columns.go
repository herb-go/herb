package columns

import (
	"github.com/herb-go/herb/model/sql/db"
)

var drivers = map[string]func() ColumnsLoader{}

type Column struct {
	Field      string
	ColumnType string
	AutoValue  bool
	PrimayKey  bool
	NotNull    bool
}

type ColumnsLoader interface {
	Columns() ([]*Column, error)
	Load(conn db.Database, table string) error
}

func Register(name string, loader func() ColumnsLoader) {
	drivers[name] = loader
}

func Driver(name string) func() ColumnsLoader {
	return drivers[name]
}
