package columns

import (
	"fmt"

	"github.com/herb-go/herb/model/sql/db"
)

var drivers = map[string]func() Loader{}

//Column sql column info
type Column struct {
	//Raw raw column name
	Raw string
	//Field column name
	Field string
	//ColumnType golang type mapped to column data
	ColumnType string
	//AutoValue  if value is auto created
	AutoValue bool
	//PrimayKey if column is primay key
	PrimayKey bool
	//NotNull if column is not null column
	NotNull bool
}

//Loader loader which load columns info
type Loader interface {
	// Columns return loaded columns
	Columns() ([]*Column, error)
	// Load load columns with given database and table name
	Load(conn db.Database, table string, fieldPrefix string) error
}

//Register register columns loader with given name
func Register(name string, loader func() Loader) {
	drivers[name] = loader
}

//Driver get driver with given name
func Driver(name string) (Loader, error) {
	d := drivers[name]
	if d == nil {
		return nil, fmt.Errorf("column driver \"%s\" not registered", name)
	}
	return d(), nil
}
