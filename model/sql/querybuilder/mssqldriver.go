package querybuilder

type MSSQLBuilderDriver struct {
	EmptyBuilderDriver
}

//LimitCommandBuilder build limit command with given limit query.
func (d *MSSQLBuilderDriver) LimitCommandBuilder(q *LimitQuery) string {
	var command = ""
	if q.offset != nil {
		command += " OFFSET ? ROWS "
	}
	if q.limit != nil {
		command += " FETCH NEXT ? ROWS ONLY "
	}

	return command
}

//LimitArgBuilder build limit args with given limit query.
func (d *MSSQLBuilderDriver) LimitArgBuilder(q *LimitQuery) []interface{} {
	var args = []interface{}{}
	if q.limit != nil {
		args = append(args, *q.limit)
	}
	if q.offset != nil {
		args = append(args, *q.offset)
	}
	return args
}

var MSSQLDriver = &MSSQLBuilderDriver{}

func init() {
	RegisterDriver("mssql", MSSQLDriver)
}
