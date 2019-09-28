package mysql

import (
	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/herb-go/herb/model/sql/querybuilder"
)

// BuilderDriver mysql bilder driver struct
type BuilderDriver struct {
	querybuilder.EmptyBuilderDriver
}

//IsDuplicate check if error is Is duplicate error.
func (d *BuilderDriver) IsDuplicate(err error) bool {
	if err == nil {
		return false
	}
	e, ok := err.(*mysqldriver.MySQLError)
	if ok == false {
		return false
	}
	return e.Number == 1062
}
func init() {
	querybuilder.RegisterDriver("mysql", &BuilderDriver{})
}
