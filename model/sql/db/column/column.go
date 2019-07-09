package column

import (
	"github.com/herb-go/herb/model/sql/db"
)

var Drivers = map[string]func() ColumnsLoader{}

type Column struct {
	Field      string
	ColumnType string
	AutoValue  bool
	PrimayKey  bool
	NotNull    bool
}

type Columns []Column

type ColumnsLoader interface {
	Columns() ([]Column, error)
	Load(conn db.Database, table string) error
}

func Register(name string, loader func() ColumnsLoader) {
	Drivers[name] = loader
}
