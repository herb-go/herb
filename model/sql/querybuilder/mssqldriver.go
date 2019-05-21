package querybuilder

// MSSQLBuilderDriver mssql bilder driver struct
type MSSQLBuilderDriver struct {
	EmptyBuilderDriver
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
