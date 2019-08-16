package querybuilder

import driver "github.com/denisenkom/go-mssqldb"

// MSSQLBuilderDriver mssql bilder driver struct
type MSSQLBuilderDriver struct {
	EmptyBuilderDriver
}

//IsDuplicate check if error is Is duplicate error.
func (d *MSSQLBuilderDriver) IsDuplicate(err error) bool {
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
func (d *MSSQLBuilderDriver) LimitCommandBuilder(q *LimitClause) string {
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
func (d *MSSQLBuilderDriver) LimitArgBuilder(q *LimitClause) []interface{} {
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
var MSSQLDriver = &MSSQLBuilderDriver{}

func init() {
	RegisterDriver("mssql", MSSQLDriver)
}
