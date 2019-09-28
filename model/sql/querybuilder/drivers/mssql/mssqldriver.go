package mssql

import driver "github.com/denisenkom/go-mssqldb"
import "github.com/herb-go/herb/model/sql/querybuilder"

// BuilderDriver mssql bilder driver struct
type BuilderDriver struct {
	querybuilder.EmptyBuilderDriver
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
	return e.Number == 2627
}

//LimitCommandBuilder build limit command with given limit query.
func (d *BuilderDriver) LimitCommandBuilder(q *querybuilder.LimitClause) string {
	var command = ""
	if q.Offset != nil {
		command += " OFFSET ? ROWS "
	}
	if q.Limit != nil {
		command += " FETCH NEXT ? ROWS ONLY "
	}

	return command
}

//LimitArgBuilder build limit args with given limit query.
func (d *BuilderDriver) LimitArgBuilder(q *querybuilder.LimitClause) []interface{} {
	var args = []interface{}{}
	if q.Limit != nil {
		args = append(args, *q.Limit)
	}
	if q.Offset != nil {
		args = append(args, *q.Offset)
	}
	return args
}

// MSSQLDriver mssql builder driver
var MSSQLDriver = &BuilderDriver{}

func init() {
	querybuilder.RegisterDriver("mssql", MSSQLDriver)
}
