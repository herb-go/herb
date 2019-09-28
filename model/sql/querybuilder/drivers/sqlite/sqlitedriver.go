package sqlite

import "github.com/herb-go/herb/model/sql/querybuilder"
import driver "github.com/mattn/go-sqlite3"

//BuilderDriver sqlite builder driver
type BuilderDriver struct {
	querybuilder.EmptyBuilderDriver
}

//TruncateTableCommandBuilder return truncate table query.
func (d *BuilderDriver) TruncateTableCommandBuilder(t string) string {
	return "DELETE FROM " + t
}
func init() {
	querybuilder.RegisterDriver("sqlite3", &BuilderDriver{})
}

//IsDuplicate check if error is Is duplicate error.
func (d *BuilderDriver) IsDuplicate(err error) bool {
	if err == nil {
		return false
	}
	e, ok := err.(driver.Error)
	if ok == false {
		return false
	}
	return e.Code == 19
}
