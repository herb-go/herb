package querybuilder

type SqliteBuilderDriver struct {
	EmptyBuilderDriver
}

func (d *SqliteBuilderDriver) DeleteCommandBuilder(q *DeleteQuery) string {
	var command = "DELETE"
	p := q.Prefix.QueryCommand()
	if p != "" {
		command += " " + p
	}
	command += " FROM " + q.TableName
	if q.alias != "" {
		command += " AS " + q.alias
	}
	return command
}

var SqliteDriver = &SqliteBuilderDriver{}

func init() {
	RegisterDriver("sqlite3", SqliteDriver)
}
